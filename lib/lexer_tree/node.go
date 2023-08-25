package lexer_tree

import (
	"github.com/lo1kx/math_parser/lib"
	"strconv"
)

type Node struct {
	Operation lib.ExpressionToken
	Data      string
	Left      *Node
	Right     *Node
}

func NewNode(operation lib.ExpressionToken, data string) *Node {
	return &Node{
		Operation: operation,
		Data:      data,
	}
}

func (n *Node) Eval() (int, error) {
	if n.Operation == lib.CONST {
		return strconv.Atoi(n.Data)
	}

	left, err := n.Left.Eval()

	if err != nil {
		return 0, err
	}

	right, err := n.Right.Eval()

	if err != nil {
		return 0, err
	}

	return n.calculate(left, right), nil
}

func (n *Node) calculate(left int, right int) int {
	switch n.Operation {
	case lib.PLUS:
		return left + right
	case lib.MINUS:
		return left - right
	case lib.DIV:
		return left / right
	case lib.MULT:
		return left * right
	default:
		return 0
	}
}
