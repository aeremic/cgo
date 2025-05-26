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

	env := value.NewEnvironment()

	return Eval(program, env)
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

func testNullValueWrapper(t *testing.T, v value.Wrapper) bool {
	if v != NULL {
		t.Errorf("value is not NULL. Got %T (%+v)", v, v)
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

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)
		if ok {
			testIntegerValueWrapper(t, evaluated, int64(integer))
		} else {
			testNullValueWrapper(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`if (true) {
				if (true) {
					if (true) {
						return 20;
					}

					return 10;
				}

				return 1;
			}
			`, 20,
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerValueWrapper(t, evaluated, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{

		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
				}
			`, "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		errorWrapped, ok := evaluated.(*value.Error)
		if !ok {
			t.Errorf("No error returned. Got %T(%+v)", evaluated, evaluated)
			continue
		}

		if errorWrapped.Message != test.expectedMessage {
			t.Errorf("Invalid message. Got %s instead of %s",
				errorWrapped.Message, test.expectedMessage)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerValueWrapper(t, evaluated, test.expected)
	}
}
