package styledjsx

import (
	"fmt"
	"strconv"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/matthewmueller/jsx"
	"github.com/matthewmueller/jsx/ast"
	"github.com/matthewmueller/styledjsx/internal/style"
)

func New() *Rewriter {
	return &Rewriter{
		Import: Import{
			Name:      "Style",
			Path:      "styled-jsx",
			IsDefault: true,
		},
		Minify: false,
	}
}

func Rewrite(path, code string) (string, error) {
	return New().Rewrite(path, code)
}

type Import struct {
	Name      string
	Path      string
	IsDefault bool
}

func (i *Import) String() string {
	if i.IsDefault {
		return "import " + i.Name + " from " + strconv.Quote(i.Path) + ";\n"
	}
	return "import {" + i.Name + "} from " + strconv.Quote(i.Path) + ";\n"
}

type Rewriter struct {
	Import Import
	// Used by ESBuild
	Minify  bool
	Engines []esbuild.Engine
}

func (r *Rewriter) Rewrite(path, code string) (string, error) {
	ast, err := jsx.Parse(path, code)
	if err != nil {
		return "", fmt.Errorf("styledjsx: unable to parse %q: %w", path, err)
	}
	if err := r.RewriteAST(path, ast); err != nil {
		return "", fmt.Errorf("styledjsx: unable to rewrite %q: %w", path, err)
	}
	return ast.String(), nil
}

func (r *Rewriter) RewriteAST(path string, script *ast.Script) error {
	visitor := &style.Visitor{
		Path:       path,
		Prefix:     "jsx-",
		Attr:       "class",
		ImportName: r.Import.Name,
		Minify:     r.Minify,
		Engines:    r.Engines,
	}
	if err := visitor.Rewrite(script); err != nil {
		return err
	}
	// Add the import statement to the top of the file
	script.Body = append([]ast.Fragment{
		&ast.Text{Value: r.Import.String()}},
		script.Body...,
	)
	return nil
}
