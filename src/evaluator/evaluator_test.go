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

func testBooleanValueWrapper(t *testing.T, v value.Wrapper, expected bool) bool {
	result, ok := v.(*value.Boolean)
	if !ok {
		t.Errorf("v is not Boolean type. Got %T", v)

		return false
	}

	if result.Value != expected {
		t.Errorf("v has wrong value. Got %t instead of %t",
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
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerValueWrapper(t, evaluated, test.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanValueWrapper(t, evaluated, test.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!!5", false},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanValueWrapper(t, evaluated, test.expected)
	}
}
