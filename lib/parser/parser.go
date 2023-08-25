package parser

import (
	"github.com/lo1kx/math_parser/lib"
	"github.com/lo1kx/math_parser/lib/lexer_tree"
	"github.com/lo1kx/math_parser/lib/stack"
	"unicode"

	"github.com/lo1kx/math_parser/lib/validation"
)

var (
	priority = map[rune]int{
		')': 0,
		'(': 0,
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}
	charToToken = map[rune]lib.ExpressionToken{
		'+': lib.PLUS,
		'-': lib.MINUS,
		'*': lib.MULT,
		'/': lib.DIV,
	}
)

type element struct {
	Token lib.ExpressionToken
	Data  string
}

type ExpressionParser struct {
	Expression string
	CurPos     int
	Root       *lexer_tree.Node
}

func NewExpressionParser(expression string) *ExpressionParser {
	expression = validation.Purify(expression)
	return &ExpressionParser{
		Expression: expression,
		CurPos:     0,
		Root:       nil,
	}
}

func (ep *ExpressionParser) Validate() error {
	return validation.Validate(ep.Expression)
}

func (ep *ExpressionParser) Parse() error {
	err := ep.Validate()
	if err != nil {
		return err
	}

	ep.do()

	return nil
}

func (ep *ExpressionParser) do() {
	stackNode := stack.NewStack[*lexer_tree.Node]()
	stackChar := stack.NewStack[rune]()

	var (
		root *lexer_tree.Node
		tmp1 *lexer_tree.Node
		tmp2 *lexer_tree.Node
	)

	for state := ep.iterate(); state.Token != lib.EOE; {
		switch state.Token {
		case lib.PLUS, lib.MINUS, lib.DIV, lib.MULT:
			ep.operation(stackChar, stackNode, root, tmp1, tmp2, state)
		case lib.CONST:
			root = lexer_tree.NewNode(lib.CONST, state.Data)
			stackNode.Push(root)
		case lib.LEFT_BRACE:
			stackChar.Push('(')
		case lib.RIGHT_BRACE:
			ep.expressionUnwinding(stackChar, stackNode, root, tmp1, tmp2)
		}

		state = ep.iterate()
	}

	ep.expressionUnwinding(stackChar, stackNode, root, tmp1, tmp2)

	ep.Root = stackNode.Top()
}

func (ep *ExpressionParser) expressionUnwinding(
	stackChar *stack.Stack[rune],
	stackNode *stack.Stack[*lexer_tree.Node],
	root *lexer_tree.Node,
	tmp1 *lexer_tree.Node,
	tmp2 *lexer_tree.Node,
) {
	for !stackChar.Empty() && stackChar.Top() != '(' {
		root = lexer_tree.NewNode(charToToken[stackChar.Top()], string(stackChar.Top()))
		stackChar.Pop()

		tmp1 = stackNode.Top()
		stackNode.Pop()

		tmp2 = stackNode.Top()
		stackNode.Pop()

		root.Left = tmp2
		root.Right = tmp1

		stackNode.Push(root)
	}

	stackChar.Pop()
}

func (ep *ExpressionParser) operation(
	stackChar *stack.Stack[rune],
	stackNode *stack.Stack[*lexer_tree.Node],
	root *lexer_tree.Node,
	tmp1 *lexer_tree.Node,
	tmp2 *lexer_tree.Node,
	state element,
) {
	for !stackChar.Empty() && stackChar.Top() != '(' &&
		priority[stackChar.Top()] >= priority[rune(state.Data[0])] {

		root = lexer_tree.NewNode(charToToken[stackChar.Top()], string(stackChar.Top()))
		stackChar.Pop()

		tmp1 = stackNode.Top()
		stackNode.Pop()

		tmp2 = stackNode.Top()
		stackNode.Pop()

		root.Left = tmp2
		root.Right = tmp1

		stackNode.Push(root)
	}

	stackChar.Push(rune(state.Data[0]))
}

func (ep *ExpressionParser) tokenize() lib.ExpressionToken {
	if ep.CurPos >= len(ep.Expression) {
		return lib.EOE
	}

	switch ep.Expression[ep.CurPos] {
	case '(':
		return lib.LEFT_BRACE
	case ')':
		return lib.RIGHT_BRACE
	case '/':
		return lib.DIV
	case '-':
		return lib.MINUS
	case '+':
		return lib.PLUS
	case '*':
		return lib.MULT
	default:
		if unicode.IsDigit(rune(ep.Expression[ep.CurPos])) {
			return lib.CONST
		}
		return lib.INVALID
	}
}

func (ep *ExpressionParser) getNum() string {
	data := ""

	for token := ep.tokenize(); token == lib.CONST; {
		data += string(ep.Expression[ep.CurPos])
		ep.CurPos++
		token = ep.tokenize()
	}

	return data
}

func (ep *ExpressionParser) iterate() element {
	data := ep.getNum()
	if len(data) != 0 {
		return element{
			Token: lib.CONST,
			Data:  data,
		}
	}

	token := ep.tokenize()
	if token == lib.EOE {
		return element{
			Token: token,
			Data:  "",
		}
	}

	ep.CurPos++
	return element{
		Token: token,
		Data:  string(ep.Expression[ep.CurPos-1]),
	}
}
