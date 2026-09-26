package spreadsheet

import (
	"errors"
	"slices"
)

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

func (a *Area) Values() []int {
	// useful for developing other functions working on the values (not-their positions) in some area
	vals := make([]int, 0, a.height()*a.width())
	for i := a.r1; i <= a.r2; i++ {
		for j := a.c1; j <= a.c2; j++ {
			if a.sh.grid[i][j].isValid() {
				vals = append(vals, a.sh.grid[i][j].value)
			}
		}
	}
	return vals
}

func (a *Area) Count() int {
	vals := a.Values()
	return len(vals)
}

func (a *Area) Max() (int, error) {
	vals := a.Values()
	if len(vals) == 0 {
		return 0, errors.New("There is no value in this area!")
	}
	return slices.Max((vals)), nil
}

func (a *Area) Min() (int, error) {
	vals := a.Values()
	if len(vals) == 0 {
		return 0, errors.New("There is no value in this area!")
	}
	return slices.Min((vals)), nil
}

func (a *Area) Sum() (int, error) {
	vals := a.Values()
	if len(vals) == 0 {
		return 0, errors.New("There is no value in this area!")
	}
	sum := 0
	for _, v := range vals {
		sum += v
	}
	return sum, nil
}

func (a *Area) Average() (float64, error) {
	vals := a.Values()
	if len(vals) == 0 {
		return 0, errors.New("There is no value in this area!")
	}
	sum := 0
	for _, v := range vals {
		sum += v
	}
	return float64(sum / len(vals)), nil
}
