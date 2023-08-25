package lib

type ExpressionToken int

const (
	MULT        ExpressionToken = iota
	PLUS        ExpressionToken = iota
	MINUS       ExpressionToken = iota
	DIV         ExpressionToken = iota
	CONST       ExpressionToken = iota
	LEFT_BRACE  ExpressionToken = iota
	RIGHT_BRACE ExpressionToken = iota
	INVALID     ExpressionToken = iota
	EOE         ExpressionToken = iota
)
