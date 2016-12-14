package runner

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/noonien/techon/lexer"
	"github.com/noonien/techon/parser"
)

type Machine struct {
	Addresses map[string]int
	Variables []*Variable
	Functions map[string]*parser.FunctionStatement
	Stack     []int
}

type Variable struct {
	Name string
	Size int
	Data []int
}

func NewMachine() *Machine {
	return &Machine{
		Addresses: make(map[string]int),
		Functions: make(map[string]*parser.FunctionStatement),
	}
}

func (m *Machine) Execute(st parser.Statement) error {
	return m.exec(st)
}

func (m *Machine) exec(st parser.Statement) error {
	switch st := st.(type) {
	case parser.Program:
		for _, cst := range st {
			err := m.exec(cst)
			if err != nil {
				return err
			}
		}

	case *parser.Comment:
		err := m.debugComments(st)
		if err != nil {
			return err
		}

	case *parser.DeclarationStatement:
		err := m.declareVariable(st)
		if err != nil {
			return err
		}

	case *parser.FunctionStatement:
		err := m.function(st)
		if err != nil {
			return err
		}

	case *parser.PushNumberStatement:
		err := m.pushNumber(st)
		if err != nil {
			return err
		}

	case *parser.IdentifierCallStatement:
		err := m.indentifierCall(st)
		if err != nil {
			return err
		}

	case parser.MathOperationStatement:
		err := m.mathOperation(st)
		if err != nil {
			return err
		}

	case *parser.DropStatement:
		err := m.drop(st)
		if err != nil {
			return err
		}

	case *parser.DupStatement:
		err := m.dup(st)
		if err != nil {
			return err
		}

	case *parser.SwapStatement:
		err := m.swap(st)
		if err != nil {
			return err
		}

	case parser.CompareOperationStatement:
		err := m.compare(st)
		if err != nil {
			return err
		}

	case *parser.GetStatement:
		err := m.get(st)
		if err != nil {
			return err
		}

	case *parser.StoreStatement:
		err := m.store(st)
		if err != nil {
			return err
		}

	case *parser.IfStatement:
		err := m._if(st)
		if err != nil {
			return err
		}

	case *parser.WhileStatement:
		err := m.while(st)
		if err != nil {
			return err
		}

	case *parser.QuitStatement:
		return nil

	default:
	}

	return nil
}

func (m *Machine) declareVariable(st *parser.DeclarationStatement) error {
	v := &Variable{
		Name: st.Name,
		Size: st.Cells,
		Data: make([]int, st.Cells),
	}

	if _, ok := m.Addresses[v.Name]; ok {
		return errors.New("cannot redeclare variable \"" + v.Name + "\"")
	}

	if _, ok := m.Functions[v.Name]; ok {
		return errors.New("cannot declare variable \"" + v.Name + "\", function already exists with that name")
	}

	var addr int
	if len(m.Variables) > 0 {
		lastVar := m.Variables[len(m.Variables)-1]
		addr = m.Addresses[lastVar.Name] + lastVar.Size
	}

	m.Addresses[v.Name] = addr
	m.Variables = append(m.Variables, v)
	return nil
}

func (m *Machine) function(st *parser.FunctionStatement) error {
	if _, ok := m.Addresses[st.Name]; ok {
		return errors.New("cannot define function \"" + st.Name + "\", variable with this name already exists")
	}

	if _, ok := m.Functions[st.Name]; ok {
		return errors.New("cannot redefine function \"" + st.Name + "\"")
	}

	m.Functions[st.Name] = st
	return nil
}

func (m *Machine) pushNumber(st *parser.PushNumberStatement) error {
	m.Stack = append(m.Stack, st.Number)
	return nil
}

func (m *Machine) indentifierCall(st *parser.IdentifierCallStatement) error {
	if addr, ok := m.Addresses[st.Identifier]; ok {
		m.Stack = append(m.Stack, addr)
		return nil
	}

	if fn, ok := m.Functions[st.Identifier]; ok {
		for _, st := range fn.Body {
			err := m.exec(st)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return errors.New("cannot resolve identifier \"" + st.Identifier + "\"")
}

func (m *Machine) resolveVariable(addr int) (*Variable, int, error) {
	var caddr int

	var v *Variable
	for _, v = range m.Variables {
		if caddr <= addr && addr < caddr+v.Size {
			break
		}
		caddr += v.Size
	}

	if v == nil || addr < caddr || addr >= caddr+v.Size {
		return nil, 0, errors.New("could not resolve address")
	}

	return v, addr - caddr, nil
}

func (m *Machine) resolveAddr(addr int) (*int, error) {
	v, idx, err := m.resolveVariable(addr)
	if err != nil {
		return nil, err
	}

	return &v.Data[idx], nil
}

func (m *Machine) mathOperation(st parser.MathOperationStatement) error {
	if len(m.Stack) < 2 {
		return errors.New("cannot perform math operation, stack does not have 2 items")
	}

	op1, op2 := m.Stack[len(m.Stack)-2], m.Stack[len(m.Stack)-1]

	var res int
	switch lexer.Token(st) {
	case lexer.Minus:
		res = op1 - op2
	case lexer.Plus:
		res = op1 + op2
	case lexer.Multiply:
		res = op1 * op2
	case lexer.Divide:
		res = op1 / op2
	case lexer.Modulus:
		res = op1 % op2
	}

	m.Stack = append(m.Stack[:len(m.Stack)-2], res)
	return nil
}

func (m *Machine) drop(st *parser.DropStatement) error {
	if len(m.Stack) < 1 {
		return errors.New("cannot drop, stack empty")
	}

	m.Stack = m.Stack[:len(m.Stack)-1]
	return nil
}

func (m *Machine) dup(st *parser.DupStatement) error {
	if len(m.Stack) < 1 {
		return errors.New("cannot dup, stack empty")
	}

	m.Stack = append(m.Stack, m.Stack[len(m.Stack)-1])
	return nil
}

func (m *Machine) swap(st *parser.SwapStatement) error {
	if len(m.Stack) < 2 {
		return errors.New("cannot perform swap operation, stack does not have 2 items")
	}

	idx1, idx2 := len(m.Stack)-2, len(m.Stack)-1
	m.Stack[idx1], m.Stack[idx2] = m.Stack[idx2], m.Stack[idx1]

	return nil
}

func (m *Machine) compare(st parser.CompareOperationStatement) error {
	if len(m.Stack) < 2 {
		return errors.New("cannot perform compare operation, stack does not have 2 items")
	}

	op1, op2 := m.Stack[len(m.Stack)-2], m.Stack[len(m.Stack)-1]

	var res bool
	switch lexer.Token(st) {
	case lexer.EQ:
		res = op1 == op2
	case lexer.LT:
		res = op1 < op2
	case lexer.GT:
		res = op1 > op2
	case lexer.LTE:
		res = op1 <= op2
	case lexer.GTE:
		res = op1 >= op2
	}

	val := 0
	if res {
		val = 1
	}

	m.Stack = append(m.Stack[:len(m.Stack)-2], val)
	return nil
}

func (m *Machine) get(st *parser.GetStatement) error {
	if len(m.Stack) < 1 {
		return errors.New("cannot perform if, stack empty")
	}

	addr := m.Stack[len(m.Stack)-1]

	ptr, err := m.resolveAddr(addr)
	if err != nil {
		return err
	}

	m.Stack = append(m.Stack[:len(m.Stack)-1], *ptr)
	return nil
}

func (m *Machine) store(st *parser.StoreStatement) error {
	if len(m.Stack) < 2 {
		return errors.New("cannot perform store operation, stack does not have 2 items")
	}

	val, addr := m.Stack[len(m.Stack)-2], m.Stack[len(m.Stack)-1]
	ptr, err := m.resolveAddr(addr)
	if err != nil {
		return err
	}
	*ptr = val

	m.Stack = m.Stack[:len(m.Stack)-2]
	return nil
}

func (m *Machine) _if(st *parser.IfStatement) error {
	if len(m.Stack) < 1 {
		return errors.New("cannot perform if, stack empty")
	}

	val := m.Stack[len(m.Stack)-1]
	m.Stack = m.Stack[:len(m.Stack)-1]

	if val != 0 {
		for _, st := range st.Body {
			err := m.exec(st)
			if err != nil {
				return err
			}
		}
	} else if len(st.ElseBody) > 0 {
		for _, st := range st.ElseBody {
			err := m.exec(st)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Machine) while(st *parser.WhileStatement) error {
	for {

		if len(m.Stack) < 1 {
			return errors.New("cannot perform while, stack empty")
		}

		val := m.Stack[len(m.Stack)-1]
		m.Stack = m.Stack[:len(m.Stack)-1]

		if val == 0 {
			break
		}

		for _, st := range st.Body {
			err := m.exec(st)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Machine) debugComments(st *parser.Comment) error {
	parts := strings.Split(st.Body, " ")
	if len(parts) < 2 || parts[0] != "debug" {
		return nil
	}

	switch parts[1] {
	case "stack":
		fmt.Fprint(os.Stderr, m.Stack, " ", strings.Join(parts[2:], " "), "\n")
	case "var":
		if len(parts) < 3 {
			return nil
		}

		addr, ok := m.Addresses[parts[2]]
		if !ok {
			return errors.New("invalid variable " + parts[2])
		}

		v, idx, err := m.resolveVariable(addr)
		if err != nil {
			return err
		}

		fmt.Fprint(os.Stderr, v.Name, " ", v.Data[idx], " ", strings.Join(parts[3:], " "), "\n")
	}

	return nil
}
