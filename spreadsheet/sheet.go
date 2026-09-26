package spreadsheet

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"
)

// Sheet represents a 2D grid of integers with a name.

type Sheet struct {
	grid [][]cell
	rows int
	cols int
	name string
}

func New(rows, cols int, name string) *Sheet {
	// New Blank Sheet with given {fixed} rows and cols and {unique} name
	// Attention: fixed size row and column grid, for now we don't support dynamic resizing
	// Handle invalid row and column access in SetValue and GetValue methods
	grid := make([][]cell, rows)
	for i := range grid {
		grid[i] = make([]cell, cols)
	}
	if name == "" {
		name = "Untitled"
	}
	return &Sheet{grid: grid, rows: rows, cols: cols, name: name}
}

func (s *Sheet) Size() (int, int) {
	return s.rows, s.cols
}

func (s *Sheet) Name() string {
	return s.name
}

// Rename to be implemented, uniqueness is important

func (s *Sheet) TotalArea() *Area {
	// generates the Area covering all the Sheet
	// we are sure this doesn't lead to Errors so we don't check conditions
	return &Area{sh: s, r1: 0, c1: 0, r2: s.rows - 1, c2: s.cols - 1}
}

func (s *Sheet) Resize(newRows, newCols int, ignoreDataLoss bool) (askAgain bool, err error) {
	var isShrink bool = newRows < s.rows || newCols < s.cols
	if isShrink && !ignoreDataLoss {
		return true, errors.New("Data Loss Error: You will be losing some data by shrinking the size of this spreadsheet")
	}
	newSh := New(newRows, newCols, s.name)
	if isShrink {
		// contract the sheet
		sourceArea, err := NewArea(s, 0, 0, newSh.rows-1, newSh.cols-1)
		if err != nil {
			return false, err // doesn't happen
		}
		sourceArea.CopyValuesTo(newSh.TotalArea())
	} else {
		// expand/extend the sheet
		targetArea, err := NewArea(newSh, 0, 0, s.rows-1, s.cols-1)
		if err != nil {
			return false, err // doesn't happen
		}
		s.TotalArea().CopyValuesTo(targetArea)
	}
	*s = *newSh // 0 or 100 rule of DB
	return false, nil
}

func (s *Sheet) SetValueAt(row, col, value int) error {
	if row < 0 || row >= s.rows || col < 0 || col >= s.cols {
		return errors.New("Invalid row or column")
	}
	s.grid[row][col].setValue(value)
	return nil
}

func (s *Sheet) GetValueAt(row, col int) (int, error) {
	if row < 0 || row >= s.rows || col < 0 || col >= s.cols {
		return 0, errors.New("Invalid row or column")
	}
	return s.grid[row][col].Value()
}

func (s *Sheet) DeleteValue(row, col int) error {
	if row < 0 || row >= s.rows || col < 0 || col >= s.cols {
		return errors.New("Invalid row or column")
	}
	s.grid[row][col].clear()
	return nil
}

func (s *Sheet) AllGridString() string {
	// returns a string containing ready-to-print format of the grid of this sheet
	// I searched for a good solution and found it on stackoverflow :)
	var builder strings.Builder
	w := tabwriter.NewWriter(&builder, 1, 1, 1, ' ', 0) // may needs config

	for i := 0; i < s.rows; i++ {
		for j := 0; j < s.cols; j++ {
			fmt.Fprintf(w, "%s\t", s.grid[i][j].toString())
		}
		fmt.Fprintln(w)
	}
	w.Flush() // what about using defer?
	return builder.String()
}
