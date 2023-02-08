package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Task struct {
	cells      [][]string
	colMapping map[string]int
	rowMapping map[string]int
}

func parseLines(lines []string) (*Task, error) {
	task := &Task{
		cells:      make([][]string, len(lines)),
		colMapping: make(map[string]int),
		rowMapping: make(map[string]int),
	}
	header := strings.Split(lines[0], ",")
	for i, v := range header {
		if i == 0 {
			continue
		}
		if v == "" {
			return nil, fmt.Errorf("missing column name in col %d", i)
		}
		if regexp.MustCompile(`\d`).MatchString(v) {
			return nil, fmt.Errorf("column name must not contain numbers: %s in col %d", v, i)
		}
		task.colMapping[v] = i

	}

	for i, line := range lines {
		splitLine := strings.Split(line, ",")
		if len(splitLine) != len(header) {
			return nil, fmt.Errorf("mismatched row lengths: row %d has length %d instead of %d", i, len(splitLine), len(header))
		}
		task.cells[i] = append(task.cells[i], splitLine...)
		task.rowMapping[task.cells[i][0]] = i

		if i != 0 {
			for j, v := range task.cells[i] {
				if v == "" {
					return nil, fmt.Errorf("empty cell at [%d][%d]", i, j)
				}
				if j != 0 {
					continue
				}
				if _, err := strconv.Atoi(v); err != nil {
					return nil, fmt.Errorf("row number must be numerical: %s", v)
				}
			}
		}
	}

	return task, nil
}

func (task *Task) getCellRawValue(address string) (string, error) {
	ptrOk := func(ptr int) bool {
		return task.colMapping[address[0:ptr]] != 0 &&
			task.rowMapping[address[ptr:]] != 0
	}

	curPtr := 1
	for curPtr < len(address) && !ptrOk(curPtr) {
		curPtr++
	}
	if !ptrOk(curPtr) {
		return "", fmt.Errorf("invalid cell address: not found: %s", address)
	}

	col := task.colMapping[address[0:curPtr]]
	row := task.rowMapping[address[curPtr:]]

	return task.cells[row][col], nil
}

func (task *Task) getOperandValue(operand string) (int, error) {
	value, err := strconv.Atoi(operand)
	if err == nil {
		return value, nil
	}

	cellRawValue, err := task.getCellRawValue(operand)
	if err != nil {
		return 0, err
	}

	value, err = strconv.Atoi(cellRawValue)
	if err == nil {
		return value, nil
	}

	exp, err := parseCellExpression(cellRawValue)
	if err != nil {
		return 0, err
	}
	if exp.left == operand || exp.right == operand {
		return 0, fmt.Errorf("invalid cell address: recursion: %s value is %s", operand, cellRawValue)
	}

	value, err = task.calcExpression(exp)
	return value, err
}

func (task *Task) calcExpression(exp Expression) (int, error) {
	left, err := task.getOperandValue(exp.left)
	if err != nil {
		return 0, err
	}
	right, err := task.getOperandValue(exp.right)
	if err != nil {
		return 0, err
	}

	var result int

	switch exp.op {
	case ADD:
		result = left + right
	case MINUS:
		result = left - right
	case MUL:
		result = left * right
	case DIV:
		if right == 0 {
			return 0, fmt.Errorf("invalid expression %v: division by zero", exp)
		}
		result = left / right
	}
	return result, nil
}

func (task *Task) processCells() error {
	for i, line := range task.cells {
		for j, cell := range line {
			if i == 0 || j == 0 {
				continue
			}
			if cell[0] != '=' {
				if _, err := strconv.Atoi(cell); err != nil {
					return fmt.Errorf("invalid cell value: %v", cell)
				}
				continue
			}

			exp, err := parseCellExpression(cell)
			if err != nil {
				return err
			}
			value, err := task.calcExpression(exp)
			if err != nil {
				return err
			}
			task.cells[i][j] = strconv.Itoa(value)
		}
	}
	return nil
}

func (task *Task) printCells() {
	for _, row := range task.cells {
		for j, cell := range row {
			fmt.Print(cell)
			if j != len(row)-1 {
				fmt.Print(",")
			}
		}
		fmt.Println()
	}
}
