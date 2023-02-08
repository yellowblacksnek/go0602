package main

import (
	"fmt"
)

type Operation int

const (
	_ Operation = iota
	ADD
	MINUS
	MUL
	DIV
)

var operations = map[uint8]Operation{
	'+': ADD,
	'-': MINUS,
	'*': MUL,
	'/': DIV,
}

type Expression struct {
	left, right string
	op          Operation
}

func parseCellExpression(cell string) (Expression, error) {
	if len(cell) < 4 { // =1+1
		return Expression{}, fmt.Errorf("invalid cell expression: %s", cell)
	}
	cell = cell[1:]

	curPtr := 0
	for curPtr < len(cell) && operations[cell[curPtr]] == 0 {
		curPtr++
	}
	if curPtr == len(cell) {
		return Expression{}, fmt.Errorf("invalid cell: operator not found: %s", cell)
	}
	op := operations[cell[curPtr]]

	left := cell[0:curPtr]
	if len(left) == 0 {
		return Expression{}, fmt.Errorf("invalid cell: no left operand: %s", cell)
	}

	right := cell[curPtr+1:]
	if len(right) == 0 {
		return Expression{}, fmt.Errorf("invalid cell: no right operand: %s", cell)
	}

	return Expression{
		left:  left,
		right: right,
		op:    op,
	}, nil
}
