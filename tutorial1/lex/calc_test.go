package lex

import (
	"testing"
)

func TestFinal(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
	}
	for _, tt := range tests {
		res := calc(tt.input)
		if res != tt.expected {
			t.Errorf("Wrong answer, got=%d, want=%d", res, tt.expected)
		}
	}
}

func TestTokenizer(t *testing.T) {

}
