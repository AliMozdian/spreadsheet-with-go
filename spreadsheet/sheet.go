package spreadsheet

import "errors"

// Cell represents a cell in the sheet that can hold an integer value.

type cell struct {
	value int
	valid bool
}

func (c cell) IsValid() bool {
	return c.valid
}

func (c cell) Value() (int, error) {
	if !c.valid {
		return 0, errors.New("Cell is not valid")
	}
	return c.value, nil
}

func (c *cell) SetValue(value int) {
	c.value = value
	c.valid = true
}

func (c *cell) Clear() {
	c.valid = false
}

// Sheet represents a 2D grid of integers with a name.

type Sheet struct {
	grid [][]cell
	rows int
	cols int
	Name string
}

func NewBlankSheet(rows, cols int, name string) *Sheet {
	// fixed size row and column grid, for now we don't support dynamic resizing
	// handle invalid row and column access in SetValue and GetValue methods
	grid := make([][]cell, rows)
	for i := range grid {
		grid[i] = make([]cell, cols)
	}
	if name == "" {
		name = "Untitled"
	}
	return &Sheet{grid: grid, rows: rows, cols: cols, Name: name}
}

func (s *Sheet) SetValueAt(row, col, value int) error {
	if row < 0 || row >= s.rows || col < 0 || col >= s.cols {
		return errors.New("Invalid row or column")
	}
	s.grid[row][col].SetValue(value)
	return nil
}

func (s *Sheet) GetValueAt(row, col int) (int, error) {
	if row < 0 || row >= s.rows || col < 0 || col >= s.cols {
		return 0, errors.New("Invalid row or column")
	}
	return s.grid[row][col].Value()
}

func (s *Sheet) CountValidsIn(r1, r2, c1, c2 int) int {
	count := 0
	for i := r1; i <= r2; i++ {
		for j := c1; j <= c2; j++ {
			if s.grid[i][j].IsValid() {
				count++
			}
		}
	}
	return count
}
