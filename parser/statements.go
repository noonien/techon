package parser

import "github.com/noonien/techon/lexer"

type Program []Statement

type Statement interface {
}

type DeclarationStatement struct {
	Name  string
	Cells int
}

type PushNumberStatement struct {
	Number int
}

type IdentifierCallStatement struct {
	Identifier string
}

type DropStatement struct{}

type DupStatement struct{}

type SwapStatement struct{}

type Comment struct {
	Body string
}

type GetStatement struct{}

type StoreStatement struct{}

type MathOperationStatement lexer.Token

type CompareOperationStatement lexer.Token

type FunctionStatement struct {
	Name string
	Body []Statement
}

type IfStatement struct {
	Body     []Statement
	ElseBody []Statement
}

type WhileStatement struct {
	Body []Statement
}

type QuitStatement struct{}
