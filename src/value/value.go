package value

import "fmt"

type Type string

const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	NULL    = "NULL"
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
