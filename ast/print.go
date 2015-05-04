package ast

// Print AST
type PrintVisitor struct {
	buf string

	indent int
}

func NewPrintVisitor() *PrintVisitor {
	return &PrintVisitor{}
}

func PrintNode(node Node) string {
	visitor := NewPrintVisitor()
	node.Accept(visitor)
	return visitor.Output()
}

func (v *PrintVisitor) Output() string {
	return v.buf
}

func (v *PrintVisitor) add(val string) {
	for i := 0; i < v.indent; {
		v.buf += " "
		i++
	}

	v.buf += val

	v.buf += "\n"
}

//
// Visitor interface
//

func (v *PrintVisitor) visitProgram(node *ProgramNode) {
	for _, n := range node.Statements {
		n.Accept(v)
	}
}

func (v *PrintVisitor) visitContent(node *ContentNode) {
	v.add("CONTENT[" + node.Content + "]")
}

func (v *PrintVisitor) visitComment(node *CommentNode) {
	v.add("{{! '" + node.Comment + "' }}")
}
