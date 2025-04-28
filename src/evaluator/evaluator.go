package evaluator

import (
	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/value"
)

// Only once created when referenced.
var (
	TRUE  = &value.Boolean{Value: true}
	FALSE = &value.Boolean{Value: false}
)

func evalStatements(statements []ast.Statement) value.Wrapper {
	var result value.Wrapper

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

func Eval(node ast.Node) value.Wrapper {
	switch node := node.(type) {
	case *ast.ProgramRoot:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &value.Integer{
			Value: node.Value,
		}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		} else {
			return FALSE
		}
	}

	return nil
}
