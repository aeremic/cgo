package value

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/aeremic/cgo/ast"
)

type Type string
type BuiltInFunction func(args ...Wrapper) Wrapper

const (
	INTEGER  = "INTEGER"
	STRING   = "STRING"
	BOOLEAN  = "BOOLEAN"
	NULL     = "NULL"
	RETURN   = "RETURN"
	ERROR    = "ERROR"
	FUNCTION = "FUNCTION"
	BUILTIN  = "BUILTIN"
	ARRAY    = "ARRAY"
	DICT     = "DICT"
)

type Wrapper interface {
	Type() Type
	Sprintf() string
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  Type
	Value uint64
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER
}

func (i *Integer) Sprintf() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING
}

func (s *String) Sprintf() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	hash := fnv.New64()
	hash.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: hash.Sum64()}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN
}

func (b *Boolean) Sprintf() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Type() Type {
	return NULL
}

func (n *Null) Sprintf() string {
	return "null"
}

type ReturnValue struct {
	Value Wrapper
}

func (rv *ReturnValue) Type() Type {
	return RETURN
}

func (rv *ReturnValue) Sprintf() string {
	return rv.Value.Sprintf()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ERROR
}

func (e *Error) Sprintf() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type {
	return FUNCTION
}

func (f *Function) Sprintf() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BuiltIn struct {
	Fn BuiltInFunction
}

func (bi *BuiltIn) Type() Type {
	return BUILTIN
}

func (bi *BuiltIn) Sprintf() string {
	return "builtin function"
}

type Array struct {
	Elements []Wrapper
}

func (a *Array) Type() Type {
	return ARRAY
}

func (a *Array) Sprintf() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range a.Elements {
		elements = append(elements, element.Sprintf())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type DictElement struct {
	Key   Wrapper
	Value Wrapper
}

type Dict struct {
	Elements map[HashKey]DictElement
}

func (d *Dict) Type() Type {
	return DICT
}

func (d *Dict) Sprintf() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range d.Elements {
		elements = append(elements, fmt.Sprintf("%s: %s",
			element.Key.Sprintf(), element.Value.Sprintf()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}
