package lexer

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var eof = rune(0)

type Scanner struct {
	r *bufio.Reader

	line int
	col  int
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	// consume all contigous whitespace
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	}

	if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	if ch == '-' {
		next := s.read()
		s.unread()

		if isDigit(next) {
			s.unread()
			return s.scanNumber()
		}

		return Minus, string(ch)
	}

	if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	if ch == '(' {
		s.unread()
		return s.scanComment()
	}

	switch ch {
	case eof:
		return EOF, ""
	case '+':
		return Plus, string(ch)
	case '*':
		return Multiply, string(ch)
	case '/':
		return Divide, string(ch)
	case ':':
		return StartFunc, string(ch)
	case ';':
		return EndFunc, string(ch)
	case '@':
		return Get, string(ch)
	case '!':
		return Store, string(ch)

	case '<', '>', '=':
		s.unread()
		return s.scanComparator()
	}

	return ILLEGAL, string(ch)
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "VARIABLE":
		return Variable, buf.String()
	case "CELLS":
		return Cells, buf.String()
	case "IF":
		return If, buf.String()
	case "ELSE":
		return Else, buf.String()
	case "THEN":
		return Then, buf.String()
	case "WHILE":
		return While, buf.String()
	case "REPEAT":
		return Repeat, buf.String()
	case "QUIT":
		return Quit, buf.String()
	case "MOD":
		return Modulus, buf.String()
	case "DROP":
		return Drop, buf.String()
	case "DUP":
		return Dup, buf.String()
	case "SWAP":
		return Swap, buf.String()
	}

	// Otherwise return as a regular identifier.
	return Ident, buf.String()
}

// scanNumber consumes the current rune and all contiguous number runes.
func (s *Scanner) scanNumber() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	_, _ = buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Otherwise return as a regular identifier.
	return Number, buf.String()
}

// scanComparator consumes the current rune and all contiguous comparator runes.
func (s *Scanner) scanComparator() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	ch := s.read()
	_, _ = buf.WriteRune(ch)
	if ch == '=' {
		return EQ, buf.String()
	}

	next := s.read()
	if next != '=' {
		s.unread()

		switch ch {
		case '<':
			return LT, buf.String()
		case '>':
			return GT, buf.String()
		}
	}

	_, _ = buf.WriteRune(next)
	switch ch {
	case '<':
		return LTE, buf.String()
	case '>':
		return GTE, buf.String()
	}

	// Otherwise return as a regular identifier.
	return ILLEGAL, buf.String()
}

// scanComment consumes the current rune and all contiguous comment runes.
func (s *Scanner) scanComment() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// scan start of comment
	_, _ = buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		ch := s.read()
		if ch == eof {
			break
		}

		buf.WriteRune(ch)

		if ch == ')' {
			break
		}
	}

	// Otherwise return as a regular identifier.
	return Comment, buf.String()
}
