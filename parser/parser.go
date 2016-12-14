package parser

import (
	"errors"
	"io"
	"strconv"

	"github.com/noonien/techon/lexer"
)

type lex struct {
	tok lexer.Token
	lit string
}

type Parser struct {
	s              *lexer.Scanner
	buf            []lex
	actual, latest int
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: lexer.NewScanner(r), buf: make([]lex, 100)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (lexer.Token, string) {
	if p.actual != p.latest {
		lex := p.buf[p.actual]
		p.actual = (p.actual + 1) % len(p.buf)
		return lex.tok, lex.lit
	}

	tok, lit := p.s.Scan()
	for tok == lexer.WS {
		tok, lit = p.s.Scan()
	}

	p.buf[p.actual] = lex{tok, lit}
	p.latest = (p.latest + 1) % len(p.buf)
	p.actual = p.latest
	return tok, lit
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.actual--
	if p.actual < 0 {
		p.actual = len(p.buf) - 1
	}
}

func (p *Parser) Parse() (Program, error) {
	return p.parseProgram()
}

func (p *Parser) parseProgram() (Program, error) {
	var prog Program

	for {
		st, err := p.parseCommon()
		if err != nil {
			return nil, err
		}
		if st != nil {
			prog = append(prog, st)
			continue
		}

		tok, _ := p.scan()
		switch tok {
		case lexer.EOF:
			return prog, nil

		case lexer.Variable:
			p.unscan()
			st, err := p.parseVariableDeclaration()
			if err != nil {
				return nil, err
			}
			prog = append(prog, st)

		case lexer.StartFunc:
			p.unscan()
			st, err := p.parseFunc()
			if err != nil {
				return nil, err
			}
			prog = append(prog, st)

		default:
			return nil, errors.New("found invalid token: " + tok.String())
		}
	}
}

func (p *Parser) parseCommon() (Statement, error) {
	tok, lit := p.scan()
	switch tok {
	case lexer.Number:
		p.unscan()
		return p.parsePushNumber()

	case lexer.Ident:
		p.unscan()
		return p.parseIdentifierCall()

	case lexer.Minus, lexer.Plus, lexer.Multiply, lexer.Divide, lexer.Modulus:
		p.unscan()
		return p.parseMathOperation()

	case lexer.EQ, lexer.LT, lexer.GT, lexer.LTE, lexer.GTE:
		p.unscan()
		return p.parseCompareOperation()

	case lexer.Drop:
		return &DropStatement{}, nil

	case lexer.Dup:
		return &DupStatement{}, nil

	case lexer.Swap:
		return &SwapStatement{}, nil

	case lexer.Comment:
		return &Comment{Body: string(lit[1 : len(lit)-1])}, nil

	case lexer.Get:
		return &GetStatement{}, nil

	case lexer.Store:
		return &StoreStatement{}, nil

	case lexer.If:
		p.unscan()
		return p.parseIfStatement()

	case lexer.While:
		p.unscan()
		return p.parseWhileStatement()

	case lexer.Quit:
		return &QuitStatement{}, nil
	}

	p.unscan()
	return nil, nil
}

func (p *Parser) parseVariableDeclaration() (*DeclarationStatement, error) {
	// discard lexer.Variable
	p.scan()

	tok, lit := p.scan()
	if tok != lexer.Ident {
		return nil, errors.New("expected variable identifier")
	}

	st := &DeclarationStatement{
		Name:  lit,
		Cells: 1,
	}

	tok, nr := p.scan()
	ntok, _ := p.scan()
	if tok == lexer.Number && ntok == lexer.Cells {
		cells, err := strconv.Atoi(nr)
		if err != nil {
			return nil, err
		}

		if cells <= 0 {
			return nil, errors.New("array cannot have less than 1 cell")
		}

		st.Cells = cells
	} else {
		p.unscan()
		p.unscan()
	}

	return st, nil
}

func (p *Parser) parsePushNumber() (*PushNumberStatement, error) {
	_, lit := p.scan()
	nr, err := strconv.Atoi(lit)
	if err != nil {
		return nil, err
	}

	return &PushNumberStatement{
		Number: nr,
	}, nil
}

func (p *Parser) parseIdentifierCall() (*IdentifierCallStatement, error) {
	_, name := p.scan()

	return &IdentifierCallStatement{
		Identifier: name,
	}, nil
}

func (p *Parser) parseMathOperation() (MathOperationStatement, error) {
	tok, _ := p.scan()

	return MathOperationStatement(tok), nil
}

func (p *Parser) parseCompareOperation() (CompareOperationStatement, error) {
	tok, _ := p.scan()

	return CompareOperationStatement(tok), nil
}

func (p *Parser) parseFunc() (*FunctionStatement, error) {
	// scan FuncStart
	p.scan()

	// get function name
	tok, lit := p.scan()
	if tok != lexer.Ident {
		return nil, errors.New("expected function identifier")
	}

	fn := &FunctionStatement{
		Name: lit,
	}

	for {
		st, err := p.parseCommon()
		if err != nil {
			return nil, err
		}
		if st != nil {
			fn.Body = append(fn.Body, st)
			continue
		}

		tok, _ = p.scan()
		switch tok {
		case lexer.EndFunc:
			return fn, nil

		default:
			return nil, errors.New("found invalid token: " + tok.String())
		}
	}
}

func (p *Parser) parseIfStatement() (*IfStatement, error) {
	// scan If
	p.scan()

	ifst := &IfStatement{}

	body := &ifst.Body
	for {
		st, err := p.parseCommon()
		if err != nil {
			return nil, err
		}
		if st != nil {
			*body = append(*body, st)
			continue
		}

		tok, _ := p.scan()
		switch tok {
		case lexer.Then:
			return ifst, nil

		case lexer.Else:
			if body == &ifst.ElseBody {
				return nil, errors.New("already in else")
			}

			body = &ifst.ElseBody

		default:
			return nil, errors.New("found invalid token: " + tok.String())
		}
	}
}

func (p *Parser) parseWhileStatement() (*WhileStatement, error) {
	// scan While
	p.scan()

	whilest := &WhileStatement{}

	for {
		st, err := p.parseCommon()
		if err != nil {
			return nil, err
		}
		if st != nil {
			whilest.Body = append(whilest.Body, st)
			continue
		}

		tok, _ := p.scan()
		switch tok {
		case lexer.Repeat:
			return whilest, nil

		default:
			return nil, errors.New("found invalid token: " + tok.String())
		}
	}
}
