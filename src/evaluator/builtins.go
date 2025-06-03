package evaluator

import "github.com/aeremic/cgo/value"

var builtins = map[string]*value.BuiltIn{
	"len": {
		Fn: func(args ...value.Wrapper) value.Wrapper {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d",
					len(args), 1)
			}

			switch arg := args[0].(type) {
			case *value.String:
				return &value.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
}
