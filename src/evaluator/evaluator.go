package evaluator

import (
	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/value"
)

func Eval(node ast.Node) value.Wrapper {
	switch node := node.(type) {
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.ProgramRoot:
		return evalProgramRoot(node.Statements)
	case *ast.Boolean:
		return nativeBoolToBoolean(node.Value)
	case *ast.IntegerLiteral:
		return &value.Integer{
			Value: node.Value,
		}
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &value.ReturnValue{
			Value: val,
		}
	}

	return nil
}

func evalProgramRoot(statements []ast.Statement) value.Wrapper {
	var result value.Wrapper

	for _, statement := range statements {
		result = Eval(statement)

		if returnValue, ok := result.(*value.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}
func evalPrefixExpression(operator string, right value.Wrapper) value.Wrapper {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
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
	case "<":
		return nativeBoolToBoolean(lv < rv)
	case ">":
		return nativeBoolToBoolean(lv > rv)
	case "==":
		return nativeBoolToBoolean(lv == rv)
	case "!=":
		return nativeBoolToBoolean(lv != rv)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalInfixExpression(operator string, left value.Wrapper, right value.Wrapper) value.Wrapper {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == value.INTEGER && right.Type() == value.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBoolean(left == right)
	case operator == "!=":
		return nativeBoolToBoolean(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBlockStatements(block *ast.BlockStatement) value.Wrapper {
	var result value.Wrapper

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil && result.Type() == value.RETURN {
			return result
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression) value.Wrapper {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}
