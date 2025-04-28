package evaluator

import (
	"github.com/aeremic/cgo/parser"
	"github.com/aeremic/cgo/tokenizer"
	"github.com/aeremic/cgo/value"

	"testing"
)

func testEval(input string) value.Wrapper {
	t := tokenizer.New(input)
	p := parser.New(t)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerValueWrapper(t *testing.T, v value.Wrapper, expected int64) bool {
	result, ok := v.(*value.Integer)
	if !ok {
		t.Errorf("v is not Integer type. Got %T", v)

		return false
	}

	if result.Value != expected {
		t.Errorf("v has wrong value. Got %d instead of %d",
			result.Value, expected)

		return false
	}

	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerValueWrapper(t, evaluated, test.expected)
	}
}
