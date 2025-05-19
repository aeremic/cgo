package evaluator

import (
	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/value"
)

// Only once created. Reused when referenced again.
var (
	NULL  = &value.Null{}
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
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
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
		return NULL
	}

	v := right.(*value.Integer).Value
	return &value.Integer{
		Value: -v,
	}
}

func evalPrefixExpression(operator string, right value.Wrapper) value.Wrapper {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left value.Wrapper, right value.Wrapper) value.Wrapper {
	lv := left.(*value.Integer).Value
	rv := right.(*value.Integer).Value

	switch operator {
	case "+":
		return &value.Integer{
			Value: lv + rv,
		}
	case "-":
		return &value.Integer{
			Value: lv - rv,
		}
	case "*":
		return &value.Integer{
			Value: lv * rv,
		}
	case "/":
		return &value.Integer{
			Value: lv / rv,
		}
	default:
		return NULL
	}
}

func evalInfixExpression(operator string, left value.Wrapper, right value.Wrapper) value.Wrapper {
	switch {
	case left.Type() == value.INTEGER && right.Type() == value.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return NULL
	}
}
