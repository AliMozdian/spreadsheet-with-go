package spreadsheet

import (
	"errors"
)

// Sheet represents a 2D grid of integers with a name.

type Sheet struct {
	grid [][]cell
	rows int
	cols int
	Name string
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
	return &Sheet{grid: grid, rows: rows, cols: cols, Name: name}
}

func (s *Sheet) TotalArea() *Area {
	// generates the Area covering all the Sheet
	// we are sure this doesn't lead to Errors so we don't check conditions
	return &Area{sh: s, r1: 0, c1: 0, r2: s.rows - 1, c2: s.cols - 1}
}

func (s *Sheet) Resize(newRows, newCols int, ignoreDataLoss bool) error {
	var isShrink bool = newRows < s.rows || newCols > s.cols
	if isShrink && !ignoreDataLoss {
		return errors.New("Data Loss Error: You will be losing some data by shrinking the size of this spreadsheet")
	}
	newSh := New(newRows, newCols, s.Name)
	if isShrink {
		// contract the sheet
		sourceArea, err := NewArea(s, 0, 0, newSh.rows-1, newSh.cols-1)
		if err != nil {
			return err // doesn't happen
		}
		sourceArea.CopyValuesTo(newSh.TotalArea())
	} else {
		// expand/extend the sheet
		targetArea, err := NewArea(newSh, 0, 0, s.rows-1, s.cols-1)
		if err != nil {
			return err // doesn't happen
		}
		s.TotalArea().CopyValuesTo(targetArea)
	}
	s = newSh // 0 or 100 rule of DB
	return nil
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

func (s *Sheet) CountValidsIn(r1, r2, c1, c2 int) int {
	count := 0
	for i := r1; i <= r2; i++ {
		for j := c1; j <= c2; j++ {
			if s.grid[i][j].isValid() {
				count++
			}
		}
	}
	return count
}
