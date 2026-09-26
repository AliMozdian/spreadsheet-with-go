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

// A part of a Sheet (2D Range) -has access to the parent Sheet
type Area struct {
	sh *Sheet // the Sheet that this Area is belong to
	r1 int    // left edge
	c1 int    // top edge
	r2 int    // right edge
	c2 int    // bottom edge
}

func NewArea(sh *Sheet, r1, c1, r2, c2 int) (*Area, error) {
	// create an 2DArea from (r1, c1) and (r2, c2) and validates the indexes
	// but I'm not sure if it's the best way of handling that, something doesn't feel right!
	if r1 < 0 || r2 < r1 || r2 >= sh.rows || c2 < c1 || c1 < 0 || c2 >= sh.cols {
		return nil, errors.New("Area Arguments must follow 0 <= r1 <= r2 < rows and 0 <= c1 <= c2 < cols")
	}
	return &Area{sh: sh, r1: r1, c1: c1, r2: r2, c2: c2}, nil
}

func (a Area) height() int {
	return a.r2 - a.r1 + 1
}

func (a Area) width() int {
	return a.c2 - a.c1 + 1
}

func (a1 Area) sameShapeAs(a2 Area) bool {
	// returns true if hight and width of the two areas is the same
	// I'm not sure wether to pass by pointer is better
	// (optimization vs making sure it doesn't change it)
	return a1.height() == a2.height() && a1.width() == a2.width()
}

func (aSrc *Area) CopyValuesTo(aDst *Area) error {
	// copies the value of cells in aSrc Area, to cells in aDst Area
	if !aSrc.sameShapeAs(*aDst) {
		return errors.New("The Areas of Src & Dst of Copy should be in the same shape (equal height & width)")
	}
	// for optimization purposes I guess :))
	sSrc, sDst := aSrc.sh, aDst.sh
	h, w := aSrc.height(), aSrc.width()
	r1Src, c1Src := aSrc.r1, aSrc.c1
	r1Dst, c1Dst := aDst.r1, aDst.c1

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			sDst.grid[r1Dst+i][c1Dst+j] = sSrc.grid[r1Src+i][c1Src+j] // copies each cell
		}
	}
	return nil
}

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
