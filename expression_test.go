package main

import (
	"fmt"
	"testing"
)

func TestParseCell(t *testing.T) {
	tests := []struct {
		testString string
		shouldErr  bool
		want       Expression
	}{
		{"=A1+B1", false, Expression{"A1", "B1", ADD}},
		{"=111+B1", false, Expression{"111", "B1", ADD}},
		{"=AAA+B1", false, Expression{"AAA", "B1", ADD}},
		{"=A1+BBB", false, Expression{"A1", "BBB", ADD}},
		{"=A1-B1", false, Expression{"A1", "B1", MINUS}},
		{"=A1*B1", false, Expression{"A1", "B1", MUL}},
		{"=A1/B1", false, Expression{"A1", "B1", DIV}},
		{"=A1B1", true, Expression{}},
		{"=AAAAA", true, Expression{}},
		{"AAAAA", true, Expression{}},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s", tt.testString)
		t.Run(testName, func(t *testing.T) {
			got, err := parseCellExpression(tt.testString)
			if err != nil && !tt.shouldErr {
				t.Fatalf("parseCell(=A1+B1) returned error %s; want %v", err, tt.want)
			}
			if !tt.shouldErr && got != tt.want {
				t.Fatalf("parseCell(=A1+B1) returned %v; want %v", got, tt.want)
			}
		})
	}
}
