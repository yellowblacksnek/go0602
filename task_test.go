package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseLines(t *testing.T) {
	var testData = []string{
		",A,B,Cell",
		"1,1,0,1",
		"2,2,=A1+Cell30,0",
		"30,0,=B1+A1,5",
	}
	var want = [][]string{
		{"", "A", "B", "Cell"},
		{"1", "1", "0", "1"},
		{"2", "2", "=A1+Cell30", "0"},
		{"30", "0", "=B1+A1", "5"},
	}
	task, err := parseLines(testData)
	if err != nil {
		t.Fatalf("parseLines returned error: %s", err)
	}

	if !reflect.DeepEqual(task.cells, want) {
		t.Fatalf("parseLines got %v; want %v", task.cells, want)
	}
}

func TestParseLinesError(t *testing.T) {
	var tests = []struct {
		testName     string
		testData     []string
		errorContent string
	}{
		{testName: "missing column name", testData: []string{
			",,B,Cell",
			"1,1,0,1",
			"2,2,=A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "missing column name",
		},
		{testName: "empty row name element", testData: []string{
			",A,B,Cell",
			",1,0,1",
			"2,2,=A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "empty cell",
		},
		{testName: "empty cell element", testData: []string{
			",A,B,Cell",
			"1,1,0,1",
			"2,,=A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "empty cell"},
		{testName: "incorrect row length (less)", testData: []string{
			",A,B,Cell",
			"1,1,0",
			"2,2,=A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "mismatched row lengths"},
		{testName: "incorrect row length (more)", testData: []string{
			",A,B,Cell",
			"1,1,0,1,2",
			"2,2,=A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "mismatched row lengths"},
		{testName: "column name has digits", testData: []string{
			",A1,B,Cell",
			"1,1,0,1",
			"2,2,A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "must not contain numbers"},
		{testName: "row number has non-digits", testData: []string{
			",A,B,Cell",
			"1d,1,0,1",
			"2,2,A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "must be numerical"},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			task, err := parseLines(test.testData)
			if err == nil {
				t.Fatalf("parseLines got %v; should error", task.cells)
			}
			if !strings.Contains(err.Error(), test.errorContent) {
				t.Fatalf("parseLines error got \"%s\"; want to contain \"%s\"", err.Error(), test.errorContent)
			}
		})
	}
}

func TestGetCellRawValue(t *testing.T) {
	var testData = []string{
		",A,B,Cell",
		"1,1,0,1",
		"2,2,=A1+Cell30,0",
		"30,0,=B1+A1,5",
	}
	tests := map[string]string{
		"A1":     "1",
		"A2":     "2",
		"A30":    "0",
		"B1":     "0",
		"B2":     "=A1+Cell30",
		"B30":    "=B1+A1",
		"Cell1":  "1",
		"Cell2":  "0",
		"Cell30": "5",
	}
	task, err := parseLines(testData)
	if err != nil {
		t.Fatalf("parseLines returned error: %s", err)
	}

	for address, want := range tests {
		t.Run(address, func(t *testing.T) {
			got, err := task.getCellRawValue(address)
			if err != nil {
				t.Fatalf("getCellRawValue returned error: %s", err)
			}
			if got != want {
				t.Fatalf("getCellRawValue got %s; want %s", got, want)
			}
		})
	}
}

func TestProcessLines(t *testing.T) {
	tests := []struct {
		testName string
		testData []string
		want     [][]string
	}{
		{testName: "regular", testData: []string{
			",A,B,Cell",
			"1,1,0,1",
			"2,2,=A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, want: [][]string{
			{"", "A", "B", "Cell"},
			{"1", "1", "0", "1"},
			{"2", "2", "6", "0"},
			{"30", "0", "1", "5"},
		},
		},
		{testName: "one simple value in expression", testData: []string{
			",A,B,Cell",
			"1,1,0,1",
			"2,2,=10+Cell30,0",
			"30,0,=B1+A1,5",
		}, want: [][]string{
			{"", "A", "B", "Cell"},
			{"1", "1", "0", "1"},
			{"2", "2", "15", "0"},
			{"30", "0", "1", "5"},
		},
		}, {testName: "two simple values in expression", testData: []string{
			",A,B,Cell",
			"1,1,0,1",
			"2,2,=10+30,0",
			"30,0,=B1+A1,5",
		}, want: [][]string{
			{"", "A", "B", "Cell"},
			{"1", "1", "0", "1"},
			{"2", "2", "40", "0"},
			{"30", "0", "1", "5"},
		},
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			task, err := parseLines(test.testData)
			if err != nil {
				t.Fatalf("parseLines returned error: %s", err)
			}
			err = task.processCells()
			if err != nil {
				t.Fatalf("processCells returned error: %s", err)
			}
			got := task.cells
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("processCells got %v; want %v", got, test.want)
			}
		})
	}
}

func TestProcessLinesError(t *testing.T) {
	var tests = []struct {
		testName     string
		testData     []string
		errorContent string
	}{
		{testName: "incorrect cell", testData: []string{
			",A,B,Cell",
			"1,1,0,1",
			"2,2,A1+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "invalid cell value",
		},
		{testName: "recursion", testData: []string{
			",A,B,Cell",
			"1,=A1+B1,0,1",
			"2,2,7,0",
			"30,0,9,5",
		}, errorContent: "recursion",
		},
		{testName: "non existing cell address", testData: []string{
			",A,B,Cell",
			"1,1,0,1",
			"2,2,=A5+Cell30,0",
			"30,0,=B1+A1,5",
		}, errorContent: "not found"},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			task, err := parseLines(test.testData)
			if err != nil {
				t.Fatalf("parseLines got error %s", err.Error())
			}
			err = task.processCells()
			if err == nil {
				t.Fatalf("parseLines got %v; should error", task.cells)
			}
			if !strings.Contains(err.Error(), test.errorContent) {
				t.Fatalf("processCells error got \"%s\"; want to contain \"%s\"", err.Error(), test.errorContent)
			}
		})
	}
}
