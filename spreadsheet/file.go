package spreadsheet

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func (s *Sheet) SaveToFile(filename string) error {
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
					return errors.New("Error While Writing Values: " + err.Error())
				}
			} else {
				_, err := file.WriteString("")
				if err != nil {
					return errors.New("Error While Writing Empties: " + err.Error())
				}
			}
			if j < s.cols-1 {
				_, err := file.WriteString(",")
				if err != nil {
					return errors.New("Error While Writing Commas: " + err.Error())
				}
			}
		}
		_, err := file.WriteString("\n")
		if err != nil {
			return errors.New("Error While Writing Newlines: " + err.Error())
		}
	}

	return nil
}

func LoadFromFile(filename string, name string) (*Sheet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("Error While Loading File: " + err.Error())
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
					return nil, errors.New("Error While Converting str2int: " + err.Error())
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
