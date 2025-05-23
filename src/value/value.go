package value

import "fmt"

type Type string

const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	NULL    = "NULL"
	RETURN  = "RETURN"
	ERROR   = "ERROR"
)

type Wrapper interface {
	Type() Type
	Sprintf() string
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

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN
}

func (b *Boolean) Sprintf() string {
	return fmt.Sprintf("%t", b.Value)
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
