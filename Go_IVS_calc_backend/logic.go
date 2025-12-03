package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	TokenNumber TokenType = iota
	TokenOperator
	TokenLParen // (
	TokenRParen // )
	TokenFunction
)

// Token represents a lexical unit
type Token struct {
	Type  TokenType
	Value string
}

// precedence returns the precedence of an operator
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/", "%": // modulo usually has same precedence as * /
		return 2
	case "^":
		return 3
	}
	return 0
}

// isOperator checks for a supported operator
func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/" || s == "%" || s == "^"
}

// Tokenize splits the expression into tokens
func Tokenize(expr string) ([]Token, error) {
	var tokens []Token
	runes := []rune(expr)

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		if unicode.IsSpace(char) {
			continue
		}

		if unicode.IsDigit(char) || char == '.' {
			// Parse number
			start := i
			decimal := false
			if char == '.' {
				decimal = true
			}
			for i+1 < len(runes) && (unicode.IsDigit(runes[i+1]) || runes[i+1] == '.') {
				if decimal && runes[i+1] == '.' {
					return nil, fmt.Errorf("unexpected number of decimal points")
				} else if !decimal && runes[i+1] == '.' {
					decimal = true
				}
				i++
			}
			tokens = append(tokens, Token{Type: TokenNumber, Value: string(runes[start : i+1])})
		} else if char == '(' {
			tokens = append(tokens, Token{Type: TokenLParen, Value: "("})
		} else if char == ')' {
			tokens = append(tokens, Token{Type: TokenRParen, Value: ")"})
		} else if isOperator(string(char)) {
			// Check for negative numbers: - at start or after an operator/Lparen
			if char == '-' && (len(tokens) == 0 || tokens[len(tokens)-1].Type == TokenOperator || tokens[len(tokens)-1].Type == TokenLParen) {
				// Negative number logic check: is the next char a digit?
				if i+1 < len(runes) && (unicode.IsDigit(runes[i+1]) || runes[i+1] == '.') {
					// It's part of a number, continue to parse number loop
					start := i
					i++ // consume '-'
					decimal := false
					if runes[i] == '.' {
						decimal = true
					}
					for i+1 < len(runes) && (unicode.IsDigit(runes[i+1]) || runes[i+1] == '.') {
						if decimal && runes[i+1] == '.' {
							return nil, fmt.Errorf("unexpected number of decimal points")
						} else if !decimal && runes[i+1] == '.' {
							decimal = true
						}
						i++
					}
					tokens = append(tokens, Token{Type: TokenNumber, Value: string(runes[start : i+1])})
					continue
				}
			}
			tokens = append(tokens, Token{Type: TokenOperator, Value: string(char)})
		} else if unicode.IsLetter(char) {
			// Parse function name (e.g., sqrt)
			start := i
			for i+1 < len(runes) && unicode.IsLetter(runes[i+1]) {
				i++
			}
			tokens = append(tokens, Token{Type: TokenFunction, Value: string(runes[start : i+1])})
		} else {
			return nil, fmt.Errorf("unexpected character: %c", char)
		}
	}
	return tokens, nil
}

// ShuntingYard converts infix tokens to RPN (Reverse Polish Notation)
func ShuntingYard(tokens []Token) ([]Token, error) {
	var output []Token
	var stack []Token

	for _, token := range tokens {
		switch token.Type {
		case TokenNumber:
			output = append(output, token)
		case TokenFunction:
			stack = append(stack, token)
		case TokenOperator:
			for len(stack) > 0 && stack[len(stack)-1].Type == TokenOperator {
				top := stack[len(stack)-1]
				if precedence(top.Value) >= precedence(token.Value) {
					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, token)
		case TokenLParen:
			stack = append(stack, token)
		case TokenRParen:
			foundLParen := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.Type == TokenLParen {
					foundLParen = true
					break
				}
				output = append(output, top)
			}
			if !foundLParen {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			// If token at top of stack is a function, pop it to output
			if len(stack) > 0 && stack[len(stack)-1].Type == TokenFunction {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		if top.Type == TokenLParen {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, top)
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

// EvaluateRPN evaluates an RPN token list
func EvaluateRPN(tokens []Token) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		switch token.Type {
		case TokenNumber:
			val, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, val)
		case TokenOperator:
			if len(stack) < 2 {
				return 0, fmt.Errorf("insufficient operands for operator %s", token.Value)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var res float64
			switch token.Value {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				res = a / b
			case "%":
				res = math.Mod(a, b)
			case "^":
				res = math.Pow(a, b)
			default:
				return 0, fmt.Errorf("unknown operator %s", token.Value)
			}
			stack = append(stack, res)
		case TokenFunction:
			if len(stack) < 1 {
				return 0, fmt.Errorf("insufficient operands for function %s", token.Value)
			}
			arg := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			var res float64
			switch strings.ToLower(token.Value) {
			case "sqrt":
				if arg < 0 {
					return 0, fmt.Errorf("sqrt of negative number")
				}
				res = math.Sqrt(arg)
			case "sin":
				res = math.Sin(arg)
			case "cos":
				res = math.Cos(arg)
			case "tan":
				res = math.Tan(arg)
			default:
				return 0, fmt.Errorf("unknown function %s", token.Value)
			}
			stack = append(stack, res)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return stack[0], nil
}

// ValidateExpression checks for structural errors in the token sequence
func ValidateExpression(tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}

	// Check first token
	first := tokens[0]
	if first.Type == TokenOperator || first.Type == TokenRParen {
		return fmt.Errorf("expression cannot start with '%s'", first.Value)
	}

	// Check last token
	last := tokens[len(tokens)-1]
	if last.Type == TokenOperator || last.Type == TokenFunction || last.Type == TokenLParen {
		return fmt.Errorf("expression cannot end with '%s'", last.Value)
	}

	for i := 0; i < len(tokens)-1; i++ {
		current := tokens[i]
		next := tokens[i+1]

		switch current.Type {
		case TokenNumber:
			if next.Type == TokenNumber || next.Type == TokenFunction || next.Type == TokenLParen {
				return fmt.Errorf("unexpected token '%s' after number '%s'", next.Value, current.Value)
			}
		case TokenOperator:
			if next.Type == TokenOperator || next.Type == TokenRParen {
				return fmt.Errorf("unexpected token '%s' after operator '%s'", next.Value, current.Value)
			}
		case TokenFunction:
			if next.Type == TokenOperator || next.Type == TokenRParen {
				return fmt.Errorf("unexpected token '%s' after function '%s'", next.Value, current.Value)
			}
		case TokenLParen:
			if next.Type == TokenOperator || next.Type == TokenRParen {
				return fmt.Errorf("unexpected token '%s' after '('", next.Value)
			}
		case TokenRParen:
			if next.Type == TokenNumber || next.Type == TokenFunction || next.Type == TokenLParen {
				return fmt.Errorf("unexpected token '%s' after ')'", next.Value)
			}
		}
	}

	return nil
}

// Calculate is the main entry point for evaluating an expression string
func Calculate(expr string) (float64, error) {
	tokens, err := Tokenize(expr)
	if err != nil {
		return 0, err
	}

	if err := ValidateExpression(tokens); err != nil {
		return 0, err
	}

	rpn, err := ShuntingYard(tokens)
	if err != nil {
		return 0, err
	}
	return EvaluateRPN(rpn)
}
