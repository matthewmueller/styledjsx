package style

import (
	"errors"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/matthewmueller/css"
	"github.com/matthewmueller/css/scoper"
	"github.com/matthewmueller/jsx/ast"
	"github.com/matthewmueller/styledjsx/internal/classes"
	"github.com/matthewmueller/styledjsx/internal/murmur"
)

// func Rewrite(path, component string, minify bool, script *ast.Script) (err error) {
// 	v := &Visitor{
// 		Path:       path,
// 		Prefix:     "jsx-",
// 		Attr:       "class",
// 		ImportName: component,
// 		Minify:     minify,
// 	}
// 	script.Visit(v)
// 	return v.err
// }

type Visitor struct {
	Path       string           // the path to the file
	Prefix     string           // the prefix for the scoped styles (e.g. jsx-)
	Attr       string           // the attribute to prepend the class to (e.g. class)
	ImportName string           // the name of the import (e.g. Style)
	Minify     bool             // whether to minify the css
	Engines    []esbuild.Engine // the browser engines we're targeting

	// set while running
	err error
}

// Rewrite the script with the scoped styles
func (v *Visitor) Rewrite(script *ast.Script) error {
	v.err = nil
	script.Visit(v)
	return v.err
}

var _ ast.Visitor = &Visitor{}

func (v *Visitor) VisitScript(s *ast.Script) {
	for _, fragment := range s.Body {
		fragment.Visit(v)
	}
}

func (v *Visitor) VisitText(t *ast.Text)               {}
func (v *Visitor) VisitField(f *ast.Field)             {}
func (v *Visitor) VisitStringValue(s *ast.StringValue) {}
func (v *Visitor) VisitExpr(e *ast.Expr)               {}
func (v *Visitor) VisitBoolValue(b *ast.BoolValue)     {}
func (v *Visitor) VisitComment(n *ast.Comment)         {}

func isScopedStyle(e *ast.Element) bool {
	if e.Name != "style" {
		return false
	}
	for _, attr := range e.Attrs {
		field, ok := attr.(*ast.Field)
		if !ok {
			continue
		}
		if field.Name != "scoped" && field.Name != "jsx" {
			continue
		}
		bv, ok := field.Value.(*ast.BoolValue)
		if !ok || !bv.Value {
			continue
		}
		return true
	}
	return false
}

func getText(f ast.Fragment) string {
	switch f := f.(type) {
	case *ast.Text:
		return f.Value
	case *ast.Expr:
		sb := strings.Builder{}
		for _, fragment := range f.Fragments {
			text, ok := fragment.(*ast.Text)
			if !ok {
				continue
			}
			value := strings.TrimSpace(text.Value)
			if !isString(value) {
				continue
			}
			sb.WriteString(strings.Trim(value, "`\"'"))
		}
		return sb.String()
	default:
		return ""
	}
}

func setText(f ast.Fragment, text string) {
	switch f := f.(type) {
	case *ast.Text:
		f.Value = text
	case *ast.Expr:
		f.Fragments = []ast.Fragment{
			&ast.Text{Value: "`" + strings.ReplaceAll(text, "`", "\\`") + "`"},
		}
	}
}

func (v *Visitor) VisitElement(e *ast.Element) {
	class := ""
	// Look for <style scoped> in the children
	for _, child := range e.Children {
		el, ok := child.(*ast.Element)
		if !ok {
			continue
		} else if !isScopedStyle(el) {
			continue
		}
		cls, err := v.updateScopedStyle(el)
		if err != nil {
			v.err = err
			return
		}
		class = cls
	}
	// If we didn't find a <style scoped> element, continue
	if class == "" {
		return
	}
	// Prepend the class to the class attribute
	if err := classes.Prepend(class, v.Attr, e); err != nil {
		v.err = err
		return
	}
}

func (v *Visitor) updateScopedStyle(style *ast.Element) (class string, err error) {
	text := strings.TrimSpace(getText(style.Children[0]))
	if text == "" {
		return "", nil
	}
	hash := murmur.Hash(text)
	stylesheet, err := css.Parse(v.Path, text)
	if err != nil {
		return "", err
	}
	class = v.Prefix + hash
	scoped, err := scoper.ScopeAST(v.Path, "."+class, stylesheet)
	if err != nil {
		return "", err
	}
	css := scoped.String()
	if v.Minify || len(v.Engines) > 0 {
		result := esbuild.Transform(css, esbuild.TransformOptions{
			Loader:            esbuild.LoaderCSS,
			MinifySyntax:      v.Minify,
			MinifyWhitespace:  v.Minify,
			MinifyIdentifiers: v.Minify,
			Engines:           v.Engines,
		})
		if len(result.Errors) > 0 {
			for _, e := range result.Errors {
				err = errors.Join(err, errors.New(e.Text))
			}
			return "", err
		}
		css = string(result.Code)
	}
	// update the <style scoped> element with the component name
	style.Name = v.ImportName
	setText(style.Children[0], strings.TrimSpace(css))
	return class, nil
}

func isString(value string) bool {
	if value == "" {
		return false
	}
	switch value[0] {
	case '"', '\'', '`':
		return true
	}
	return false
}
