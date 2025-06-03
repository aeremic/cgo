package evaluator

import (
	"fmt"

	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/value"
)

// Only once created. Reused when referenced again.
var (
	NULL  = &value.Null{}
	TRUE  = &value.Boolean{Value: true}
	FALSE = &value.Boolean{Value: false}
)

func newError(format string, a ...interface{}) *value.Error {
	return &value.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(v value.Wrapper) bool {
	if v != nil {
		return v.Type() == value.ERROR
	}

	return false
}

func evalBangOperatorExpression(right value.Wrapper) value.Wrapper {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right value.Wrapper) value.Wrapper {
	if right.Type() != value.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}

	v := right.(*value.Integer).Value
	return &value.Integer{
		Value: -v,
	}
}

func evalIdentifier(node *ast.Identifier, env *value.Environment) value.Wrapper {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("%s", "identifier not found: "+node.Value)
}

func nativeBoolToBoolean(input bool) *value.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func isTruthy(v value.Wrapper) bool {
	switch v {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
