package calculation_test

import (
	"calc_service/internal/calculation"
	"testing"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		expression     string
		expectedTokens []string
	}{
		{
			expression:     "3 + 5 * (2 - 8)",
			expectedTokens: []string{"3", "+", "5", "*", "(", "2", "-", "8", ")"},
		},
		{
			expression:     "10 / 2",
			expectedTokens: []string{"10", "/", "2"},
		},
	}

	for _, tc := range testCases {
		tokens := calculation.Tokenize(tc.expression)
		if len(tokens) != len(tc.expectedTokens) {
			t.Fatalf("For expression %s, expected %d tokens, got %d", tc.expression, len(tc.expectedTokens), len(tokens))
		}
		for i, token := range tokens {
			if token != tc.expectedTokens[i] {
				t.Fatalf("For expression %s, expected token %s at position %d, got %s", tc.expression, tc.expectedTokens[i], i, token)
			}
		}
	}
}

func TestInfixToPostfix(t *testing.T) {
	testCases := []struct {
		tokens          []string
		expectedPostfix []string
		expectError     bool
	}{
		{
			tokens:          []string{"3", "+", "5", "*", "(", "2", "-", "8", ")"},
			expectedPostfix: []string{"3", "5", "2", "8", "-", "*", "+"},
			expectError:     false,
		},
		{
			tokens:          []string{"10", "/", "2"},
			expectedPostfix: []string{"10", "2", "/"},
			expectError:     false,
		},
		{
			tokens:          []string{"(", "3", "+", "4"},
			expectedPostfix: nil,
			expectError:     true,
		},
	}

	for _, tc := range testCases {
		postfix, err := calculation.InfixToPostfix(tc.tokens)
		if (err != nil) != tc.expectError {
			t.Fatalf("For tokens %v, expected error: %v, got: %v", tc.tokens, tc.expectError, err)
		}
		if !tc.expectError {
			if len(postfix) != len(tc.expectedPostfix) {
				t.Fatalf("For tokens %v, expected %d postfix tokens, got %d", tc.tokens, len(tc.expectedPostfix), len(postfix))
			}
			for i, token := range postfix {
				if token != tc.expectedPostfix[i] {
					t.Fatalf("For tokens %v, expected token %s at position %d, got %s", tc.tokens, tc.expectedPostfix[i], i, token)
				}
			}
		}
	}
}

func TestEvaluatePostfix(t *testing.T) {
	testCases := []struct {
		postfix        []string
		expectedResult float64
		expectError    bool
	}{
		{
			postfix:        []string{"3", "5", "2", "8", "-", "*", "+"},
			expectedResult: -13.0,
			expectError:    false,
		},
		{
			postfix:        []string{"10", "2", "/"},
			expectedResult: 5.0,
			expectError:    false,
		},
		{
			postfix:        []string{"3", "0", "/"},
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		result, err := calculation.EvaluatePostfix(tc.postfix)
		if (err != nil) != tc.expectError {
			t.Fatalf("For postfix %v, expected error: %v, got: %v", tc.postfix, tc.expectError, err)
		}
		if !tc.expectError && result != tc.expectedResult {
			t.Fatalf("For postfix %v, expected result %f, got %f", tc.postfix, tc.expectedResult, result)
		}
	}
}
