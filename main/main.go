package main

import (
	"fmt"
	"os"
	sheet "spreadsheetCLI/spreadsheet"
	"strconv"
)

// idea for later: add a CLI package with a Commands Struct and options (+ how many arg each one accepts)
// in that way we can use a clean general CLI manager to contribute with Sheet
// and main package will be the one using both of them by connecting the funcs in call chains

func create(crnt_sheet *sheet.Sheet, args []string) {
	if len(args) < 5 {
		fmt.Println("You should follow this pattern: create <name> <rows> <cols> --options")
		os.Exit(1)
	}

	name := args[2]
	if name == "" || name[0] == '-' {
		fmt.Println("The name cannot be empty or starts with '-'")
	}
	rows, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("The number of rows must be int!")
		os.Exit(1)
	}
	cols, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println("The number of cols must be int!")
		os.Exit(1)
	}

	sh := sheet.NewBlankSheet(rows, cols, name)
	err = sh.SaveToFile()
	if err != nil {
		fmt.Println("Error while Saving File:", err.Error())
		os.Exit(1)
	}

	// check options
	for i := 5; i < len(args); i++ {
		if args[i][0] != '-' {
			fmt.Println("options should start with '-'")
			os.Exit(1)
		}
		if args[i] == "-o" || args[i] == "--online" {
			crnt_sheet = sh // the new created sheet is loaded
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Loop Interface App is not implemented yet...")
		os.Exit(1)
	}

	var current_sheet *sheet.Sheet = nil

	cmd := os.Args[1]
	switch cmd {
	case "create":
		// create <name> <rows> <cols>
		create(current_sheet, args)
	}

}
