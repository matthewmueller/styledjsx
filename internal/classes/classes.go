package classes

import (
	"strings"

	"github.com/matthewmueller/jsx/ast"
)

// Prepend adds a class to elements in the script
func Prepend(class string, attr string, el *ast.Element) error {
	ev := &classVisitor{
		class: class,
		attr:  attr,
	}
	el.Visit(ev)
	return ev.err
}

type classVisitor struct {
	class string // jsx-123
	attr  string // class
	err   error
}

func (v *classVisitor) VisitScript(n *ast.Script)           {}
func (v *classVisitor) VisitText(n *ast.Text)               {}
func (v *classVisitor) VisitField(n *ast.Field)             {}
func (v *classVisitor) VisitStringValue(n *ast.StringValue) {}
func (v *classVisitor) VisitExpr(n *ast.Expr)               {}
func (v *classVisitor) VisitBoolValue(n *ast.BoolValue)     {}
func (v *classVisitor) VisitComment(n *ast.Comment)         {}

func (v *classVisitor) VisitElement(e *ast.Element) {
	if len(e.Name) == 0 || !isLower(e.Name[0]) {
		for _, attr := range e.Attrs {
			attr.Visit(v)
		}
		for _, child := range e.Children {
			child.Visit(v)
		}
		return
	}
	hasClass := false
	for _, attr := range e.Attrs {
		if field, ok := attr.(*ast.Field); ok && isClass(field.Name) {
			v.updateClassField(field)
			hasClass = true
		}
	}
	if !hasClass {
		e.Attrs = append(e.Attrs, &ast.Field{
			Name:  v.attr,
			Value: &ast.StringValue{Value: v.class},
		})
	}
	for _, child := range e.Children {
		child.Visit(v)
	}
}

func (v *classVisitor) updateClassField(f *ast.Field) {
	switch value := f.Value.(type) {
	case *ast.StringValue:
		updateStringValue(value, v.class)
	case *ast.Expr:
		updateExpr(value, v.class)
	}
}

func updateStringValue(value *ast.StringValue, class string) {
	value.Value = class + " " + value.Value
}

func updateExpr(expr *ast.Expr, class string) {
	if len(expr.Fragments) == 0 {
		return
	}
	first := expr.Fragments[0]
	switch first := first.(type) {
	case *ast.Text:
		first.Value = updateExprString(first.Value, class)
	}

}

func updateExprString(value, class string) string {
	switch value[0] {
	case '"':
		return "\"" + class + " " + strings.TrimLeft(value, "\"")
	case '\'':
		return "'" + class + " " + strings.TrimLeft(value, "'")
	case '`':
		return "`" + class + " " + strings.TrimLeft(value, "`")
	default:
		return "`" + class + " ${" + value + "}`"
	}
}

func isClass(name string) bool {
	return name == "class" || name == "className"
}

func isLower(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}
