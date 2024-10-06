package expression

import (
	"fmt"
	"strconv"
	"strings"
)

type tokenType int

const (
	Constant tokenType = iota
	Variable
	Operator
	Parenthesis
)

type token struct {
	Type  tokenType
	Value string
}

func Parse(expression string) (Operation, error) {
	// Clean the input equation
	exp := clean(expression)

	op, err := parse(exp)
	if err != nil {
		return nil, err
	}

	return op.Simplify(), nil
}

func clean(expression string) string {
	// Remove spaces
	return strings.ReplaceAll(expression, " ", "")
}

func parse(expression string) (Operation, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return nil, err
	}

	items := make([]struct {
		token     token
		operation Operation
	}, len(tokens))

	for i := range items {
		items[i].token = tokens[i]
	}

	// Start by creating the constants, variables and parenthesis
	// We start with them because they do not depend on other operators
	for i := range items {
		switch items[i].token.Type {
		case Constant:
			value, err := strconv.ParseFloat(items[i].token.Value, 32)
			if err != nil {
				return nil, fmt.Errorf("constant could not be parsed %s", items[i].token.Value)
			}
			items[i].operation = NewConst(float32(value))
		case Variable:
			index, err := strconv.ParseInt(items[i].token.Value, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("variable index could not be parsed %s", items[i].token.Value)
			}
			items[i].operation = NewVar(uint8(index))
		case Parenthesis:
			op, err := parse(items[i].token.Value)
			if err != nil {
				return nil, err
			}
			items[i].operation = op
		case Operator:
			continue
		default:
			return nil, fmt.Errorf("unexpected token type")
		}
	}

	// Continue with the multiplications and divisions
	// Remove the constants and variables along the way
	i := 0
	for ; i < len(items); i++ {
		if items[i].token.Type == Operator {
			switch items[i].token.Value {
			case "*", "x":
				items[i].operation = NewMul(items[i-1].operation, items[i+1].operation)
			case "/":
				items[i].operation = NewDiv(items[i-1].operation, items[i+1].operation)
			case "+", "-":
				continue
			default:
				return nil, fmt.Errorf("unexpected operator %s", items[i].token.Value)
			}

			// Remove elements at i+1, then i-1
			// Order of removal is important here
			items = append(items[:i+1], items[i+2:]...)
			items = append(items[:i-1], items[i:]...)
			if i == 0 {
				return nil, fmt.Errorf("expression cannot start with an operator")
			}
			i--
		}
	}

	// Finish with the additions and substractions
	i = 0
	for ; i < len(items); i++ {
		if items[i].token.Type == Operator {
			switch items[i].token.Value {
			case "+":
				items[i].operation = NewAdd(items[i-1].operation, items[i+1].operation)
			case "-":
				items[i].operation = NewSub(items[i-1].operation, items[i+1].operation)
			case "*", "x", "/":
				continue
			default:
				return nil, fmt.Errorf("unexpected operator %s", items[i].token.Value)
			}

			// Remove elements at i+1, then i-1
			// Order of removal is important here
			items = append(items[:i+1], items[i+2:]...)
			items = append(items[:i-1], items[i:]...)
			if i == 0 {
				return nil, fmt.Errorf("expression cannot start with an operator")
			}
			i--
		}
	}

	return items[0].operation, nil
}

func tokenize(expression string) ([]token, error) {
	var tokens []token
	for index := 0; index < len(expression); index++ {
		switch expression[index] {
		case '(':
			end, err := findParenthesisEnd(expression, index)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token{
				Type:  Parenthesis,
				Value: expression[index+1 : end],
			})
			index = end
		case '{':
			end, err := findParenthesisEnd(expression, index)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token{
				Type:  Variable,
				Value: expression[index+1 : end],
			})
			index = end
		case '+', '-', '*', 'x', '/':
			tokens = append(tokens, token{
				Type:  Operator,
				Value: string(expression[index]),
			})
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			end, err := findConstantEnd(expression, index)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token{
				Type:  Constant,
				Value: expression[index : end+1],
			})
			index = end
		case ')', ']', '}':
			return nil, fmt.Errorf("unmatched parenthesis %c", expression[index])
		default:
			return nil, fmt.Errorf("unexpected character %c", expression[index])
		}
	}

	return tokens, nil
}

func findParenthesisEnd(expression string, index int) (int, error) {
	if index >= len(expression) {
		return 0, fmt.Errorf("parenthesis index out of range")
	}

	var closing byte
	opening := expression[index]
	switch opening {
	case '(':
		closing = ')'
	case '[':
		closing = ']'
	case '{':
		closing = '}'
	default:
		return 0, fmt.Errorf("start index is not parenthesis")
	}

	count := 0
	for ; index < len(expression); index++ {
		switch expression[index] {
		case opening:
			count++
		case closing:
			count--
			if count == 0 {
				return index, nil
			}
		}
	}

	return 0, fmt.Errorf("unmatched parenthesis")
}

func findConstantEnd(expression string, index int) (int, error) {
	if index >= len(expression) {
		return 0, fmt.Errorf("constant index out of range")
	}

	for ; index < len(expression); index++ {
		switch expression[index] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			continue
		default:
			return index - 1, nil
		}
	}

	return index - 1, nil
}
