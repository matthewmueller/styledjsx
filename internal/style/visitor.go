package style

import (
	"errors"
	"strconv"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/matthewmueller/css"
	"github.com/matthewmueller/css/scoper"
	"github.com/matthewmueller/jsx/ast"
	"github.com/matthewmueller/styledjsx/internal/classes"
	"github.com/matthewmueller/styledjsx/internal/murmur"
)

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

func getCSS(fragments []ast.Fragment) string {
	sb := new(strings.Builder)
	for _, fragment := range fragments {
		switch fragment := fragment.(type) {
		case *ast.Text:
			sb.WriteString(strings.TrimSpace(fragment.Value))
		case *ast.Expr:
			for _, fragment := range fragment.Fragments {
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
		}
	}
	return strings.TrimSpace(sb.String())
}

func setCSS(style *ast.Element, css string) {
	style.Children = []ast.Fragment{
		&ast.Expr{
			Fragments: []ast.Fragment{
				&ast.Text{Value: "`" + strings.ReplaceAll(strings.TrimSpace(css), "`", "\\`") + "`"},
			},
		},
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
	text := getCSS(style.Children)
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
	style.Attrs = append(style.Attrs, &ast.Field{
		Name: "id",
		Value: &ast.StringValue{
			Raw: strconv.Quote(class),
		},
	})
	setCSS(style, css)
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
