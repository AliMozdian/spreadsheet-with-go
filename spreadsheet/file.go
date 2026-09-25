package spreadsheet

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const FORMAT string = ".csv"

func (s *Sheet) SaveToFile() error {
	filename := s.Name + FORMAT
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < s.rows; i++ {
		for j := 0; j < s.cols; j++ {
			cell := s.grid[i][j]
			if cell.isValid() {
				_, err := file.WriteString(strconv.Itoa(cell.value))
				if err != nil {
					return err
				}
			} else {
				_, err := file.WriteString("")
				if err != nil {
					return err
				}
			}
			if j < s.cols-1 {
				_, err := file.WriteString(",")
				if err != nil {
					return err
				}
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadFromFile(name string) (*Sheet, error) {
	filename := name + FORMAT
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file line by line and split by commas to get the values
	var grid [][]cell
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ",")
		row := make([]cell, len(values))
		for j, value := range values {
			if value != "" {
				intValue, err := strconv.Atoi(value)
				if err != nil {
					return nil, err
				}
				row[j].setValue(intValue)
			}
		}
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Sheet{grid: grid, rows: len(grid), cols: len(grid[0]), Name: name}, nil
}
