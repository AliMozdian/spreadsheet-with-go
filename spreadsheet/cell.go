package spreadsheet

import (
	"errors"
	"strconv"
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
	c.value = 0 // for keeping memory clean and untracable :)
}

func (c *cell) toString() string {
	// only used for in-app display, file.go doesn't use this func
	if !c.valid {
		return "."
	}
	return strconv.Itoa(c.value)
}
