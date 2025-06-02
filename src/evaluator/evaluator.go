package evaluator

import (
	"github.com/aeremic/cgo/ast"
	"github.com/aeremic/cgo/value"
)

func Eval(node ast.Node, env *value.Environment) value.Wrapper {
	switch node := node.(type) {
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ProgramRoot:
		return evalProgramRoot(node.Statements, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Boolean:
		return nativeBoolToBoolean(node.Value)
	case *ast.IntegerLiteral:
		return &value.Integer{
			Value: node.Value,
		}
	case *ast.StringLiteral:
		return &value.String{
			Value: node.Value,
		}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &value.ReturnValue{
			Value: val,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return &value.Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
	}

	return nil
}

func evalProgramRoot(statements []ast.Statement, env *value.Environment) value.Wrapper {
	var result value.Wrapper

	for _, statement := range statements {
		result = Eval(statement, env)

		switch rt := result.(type) {
		case *value.ReturnValue:
			return rt.Value
		case *value.Error:
			return rt
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

func evalStringInfixExpression(operator string, left value.Wrapper, right value.Wrapper) value.Wrapper {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	lv := left.(*value.String).Value
	rv := right.(*value.String).Value

	return &value.String{
		Value: lv + rv,
	}
}

func evalInfixExpression(operator string, left value.Wrapper, right value.Wrapper) value.Wrapper {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == value.INTEGER && right.Type() == value.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == value.STRING && right.Type() == value.STRING:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBoolean(left == right)
	case operator == "!=":
		return nativeBoolToBoolean(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBlockStatements(block *ast.BlockStatement, env *value.Environment) value.Wrapper {
	var result value.Wrapper

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == value.RETURN || rt == value.ERROR {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *value.Environment) value.Wrapper {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalExpressions(exps []ast.Expression, env *value.Environment) []value.Wrapper {
	var results []value.Wrapper

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []value.Wrapper{evaluated}
		}

		results = append(results, evaluated)
	}

	return results
}

func applyFunction(fn value.Wrapper, args []value.Wrapper) value.Wrapper {
	function, ok := fn.(*value.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := createExtendedEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func createExtendedEnv(fn *value.Function, args []value.Wrapper) *value.Environment {
	env := value.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(evaluated value.Wrapper) value.Wrapper {
	if returnValue, ok := evaluated.(*value.ReturnValue); ok {
		return returnValue.Value
	}

	return evaluated
}
