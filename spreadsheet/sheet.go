package spreadsheet

import (
	"errors"
)

// Cell represents a cell in the sheet that can hold an integer value.

type cell struct {
	value int
	valid bool
}

func (c cell) isValid() bool {
	return c.valid
}

func (c cell) Value() (int, error) {
	if !c.valid {
		return 0, errors.New("Cell is not valid")
	}
	return c.value, nil
}

func (c *cell) setValue(value int) {
	c.value = value
	c.valid = true
}

func (c *cell) clear() {
	c.valid = false
}

// Useful for Range Passing
type Area struct {
	r1 int // left edge
	r2 int // right edge
	c1 int // top edge
	c2 int // bottom edge
}

func NewArea(sh Sheet, r1, r2, c1, c2 int) (*Area, error) {
	// create an 2DArea from (r1, c1) and (r2, c2) and validates the indexes
	// but I'm not sure if it's the best way of handling that, something doesn't feel right!
	if r1 < 0 || r2 < r1 || r2 >= sh.rows || c2 < c1 || c1 < 0 || c2 >= sh.cols {
		return nil, errors.New("Area Arguments must follow 0 <= r1 <= r2 < rows and 0 <= c1 <= c2 < cols")
	}
	return &Area{r1: r1, r2: r2, c1: c1, c2: c2}, nil
}

func (a Area) height() int {
	return a.r2 - a.r1 + 1
}

func (a Area) width() int {
	return a.c2 - a.c1 + 1
}

func (a1 Area) sameShapeAs(a2 Area) bool {
	return a1.height() == a2.height() && a1.width() == a2.width()
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

func (sSrc *Sheet) CopyValuesTo(sDst *Sheet, aSrc, aDst Area) error {
	// copies the value of cells in aSrc Area from sSrc Sheet, to cells in aDst Area from sDst Sheet
	if !aSrc.sameShapeAs(aDst) {
		return errors.New("The Areas of Source and Destination of Copy should be in the same shape (equal deltaRows & deltaCols)")
	}
	// for optimization purposes I guess :))
	r1Src, c1Src := aSrc.r1, aSrc.c1
	r1Dst, c1Dst := aDst.r1, aDst.c1

	for i := 0; i < aSrc.height(); i++ {
		for j := 0; j < aSrc.width(); j++ {
			sDst.grid[r1Dst+i][c1Dst+j] = sSrc.grid[r1Src+i][c1Src+j] // copies each cell
		}
	}
	return nil
}

// func (s *Sheet) Resize(newRows, newCols int) {
// 	// to do next
// }

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
