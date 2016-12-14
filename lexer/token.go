package lexer

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS

	Variable
	Ident
	Number
	Cells

	Minus
	Plus
	Multiply
	Divide
	Modulus

	EQ
	LT
	GT
	LTE
	GTE

	Drop
	Dup
	Swap
	Comment

	Get
	Store

	If
	Else
	Then
	While
	Repeat

	StartFunc
	EndFunc

	Quit
)

func (t Token) String() string {

	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case Variable:
		return "Variable"
	case Ident:
		return "Ident"
	case Number:
		return "Number"
	case Cells:
		return "Cells"
	case Minus:
		return "Minus"
	case Plus:
		return "Plus"
	case Multiply:
		return "Multiply"
	case Divide:
		return "Divide"
	case Modulus:
		return "Modulus"
	case EQ:
		return "EQ"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case Drop:
		return "Drop"
	case Dup:
		return "Dup"
	case Swap:
		return "Swap"
	case Comment:
		return "Comment"
	case Get:
		return "Get"
	case Store:
		return "Store"
	case If:
		return "If"
	case Else:
		return "Else"
	case Then:
		return "Then"
	case While:
		return "While"
	case Repeat:
		return "Repeat"
	case StartFunc:
		return "StartFunc"
	case EndFunc:
		return "EndFunc"
	case Quit:
		return "Quit"
	}

	return "Unknown"
}
