package evaluator

import (
	"fmt"

	"github.com/aeremic/cgo/value"
)

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
			case *value.Array:
				return &value.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...value.Wrapper) value.Wrapper {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d",
					len(args), 1)
			}

			if args[0].Type() != value.ARRAY {
				return newError("argument to `first` method must be an array. got=%s",
					args[0].Type())
			}

			arr := args[0].(*value.Array)
			if len(arr.Elements) == 0 {
				return NULL
			}

			return arr.Elements[0]
		},
	},
	"last": {
		Fn: func(args ...value.Wrapper) value.Wrapper {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d",
					len(args), 1)
			}

			if args[0].Type() != value.ARRAY {
				return newError("argument to `first` method must be an array. got=%s",
					args[0].Type())
			}

			arr := args[0].(*value.Array)
			if len(arr.Elements) == 0 {
				return NULL
			}

			return arr.Elements[len(arr.Elements)-1]
		},
	},
	"tail": {
		Fn: func(args ...value.Wrapper) value.Wrapper {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=%d",
					len(args), 1)
			}

			if args[0].Type() != value.ARRAY {
				return newError("argument to `first` method must be an array. got=%s",
					args[0].Type())
			}

			arr := args[0].(*value.Array)
			if len(arr.Elements) == 0 {
				return NULL
			}

			length := len(arr.Elements)
			newElements := make([]value.Wrapper, length-1, length-1)
			copy(newElements, arr.Elements[1:length])

			return &value.Array{Elements: newElements}
		},
	},
	"push": {
		Fn: func(args ...value.Wrapper) value.Wrapper {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=%d",
					len(args), 2)
			}

			if args[0].Type() != value.ARRAY {
				return newError("argument to `first` method must be an array. got=%s",
					args[0].Type())
			}

			arr := args[0].(*value.Array)

			length := len(arr.Elements)
			newElements := make([]value.Wrapper, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements = append(newElements, args[1])

			return &value.Array{Elements: newElements}
		},
	},
	"puts": {
		Fn: func(args ...value.Wrapper) value.Wrapper {
			for _, arg := range args {
				fmt.Println(arg.Sprintf())
			}

			return NULL
		},
	},
}
