package main

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected float64
		wantErr  bool
	}{
		// Basic arithmetic
		{"Addition", "2 + 3", 5, false},
		{"Subtraction", "5 - 2", 3, false},
		{"Multiplication", "4 * 3", 12, false},
		{"Division", "10 / 2", 5, false},
		{"Decimals", "2.5 + 3.5", 6, false},

		// Order of operations
		{"Precedence 1", "2 + 3 * 4", 14, false},
		{"Precedence 2", "(2 + 3) * 4", 20, false},
		{"Precedence 3", "2 + 3 * 4 ^ 2", 50, false},

		// Functions
		{"Sqrt", "sqrt(16)", 4, false},
		{"Sin", "sin(0)", 0, false},
		{"Cos", "cos(0)", 1, false},

		// Complex expressions
		{"Complex 1", "2 * (3 + 4) - 5", 9, false},
		{"Complex 2", "sqrt(9) + 2^3", 11, false},

		// Negative numbers
		{"Negative Start", "-5 + 3", -2, false},
		{"Negative Middle", "5 + -3", 2, false},
		{"Negative Parenthesis", "5 * (-3 + 2)", -5, false},
		{"Negative Power", "-2 ^ 2", 4, false}, // Parsed as (-2)^2 because -2 is tokenized as a single number
		{"Negative Function", "sqrt(4) * -1", -2, false},
		{"Double Negative", "5 - -3", 8, false},

		// Error cases - Validation
		{"Empty", "", 0, true}, // Should return error for empty input
		{"Double Operator", "2 + + 3", 0, true},
		{"Double Number", "2 3", 0, true},
		{"Starts with Operator", "* 3", 0, true},
		{"Ends with Operator", "3 +", 0, true},
		{"Empty Parens", "()", 0, true},
		{"Mismatched Parens 1", "(2 + 3", 0, true},
		{"Mismatched Parens 2", "2 + 3)", 0, true},
		{"Unknown Function", "foo(5)", 0, true},
		{"Division by Zero", "5 / 0", 0, true},
		{"Implicit Multiplication 1", "2(3)", 0, true},
		{"Implicit Multiplication 2", "(2)(3)", 0, true},
		{"Invalid Decimal", "2.3.4", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate(%q) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				// Use a small epsilon for float comparison if needed, but simple cases usually match exactly
				if diff := got - tt.expected; diff < -0.00001 || diff > 0.00001 {
					t.Errorf("Calculate(%q) = %v, want %v", tt.expr, got, tt.expected)
				}
			}
		})
	}
}
