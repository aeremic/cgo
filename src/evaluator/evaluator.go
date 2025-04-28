package evaluator

import (
	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/value"
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
	}

	return nil
}
