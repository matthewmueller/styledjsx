package style

import "github.com/matthewmueller/jsx/ast"

type Remover struct {
	err error
}

var _ ast.Visitor = &Remover{}

// Remove the script with the scoped styles
func (r *Remover) Remove(script *ast.Script) error {
	r.err = nil
	script.Visit(r)
	return r.err
}

func (r *Remover) VisitScript(s *ast.Script) {
	for _, fragment := range s.Body {
		fragment.Visit(r)
	}
}

func (r *Remover) VisitText(t *ast.Text)               {}
func (r *Remover) VisitField(f *ast.Field)             {}
func (r *Remover) VisitStringValue(s *ast.StringValue) {}
func (r *Remover) VisitExpr(e *ast.Expr)               {}
func (r *Remover) VisitBoolValue(b *ast.BoolValue)     {}
func (r *Remover) VisitComment(n *ast.Comment)         {}

func (r *Remover) VisitElement(e *ast.Element) {
	children := []ast.Fragment{}
	for _, child := range e.Children {
		if el, ok := child.(*ast.Element); ok && isScopedStyle(el) {
			continue
		}
		children = append(children, child)
	}
	e.Children = children
}
