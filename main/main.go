package main

import (
	"fmt"
	"os"
	sheet "spreadsheetCLI/spreadsheet"
	"strconv"
	"strings"
)

// global current sheet: the sheet that is open (loaded to Program Memory)
var crntSheet *sheet.Sheet = nil

// idea for later: add a CLI package with a Commands Struct and options (+ how many arg each one accepts)
// in that way we can use a clean general CLI manager to contribute with Sheet
// and main package will be the one using both of them by connecting the funcs in call chains

func checkOpenSheet() {
	// use this in sheet-required commands
	if crntSheet == nil {
		fmt.Println("There is no open spreadsheet!")
		os.Exit(1)
	}
}

func create(args []string) {
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

			crntSheet = sh // the new created sheet is loaded
		}
	}
}

func open(args []string) {
	if len(args) != 3 {
		fmt.Println("You should follow this pattern: open <name>")
		os.Exit(1)
	}

	if crntSheet != nil {
		fmt.Printf("There is already a spreadsheet open (%s), You should close it first!\n", crntSheet.Name)
		confirmClosing()
	}

	name := args[2]
	sh, err := sheet.LoadFromFile(name)
	if err != nil {
		fmt.Println("Error While Loading File:", err.Error())
		os.Exit(2)
	}

	crntSheet = sh
	fmt.Printf("Spreadsheet (%s) is open now...\n", crntSheet.Name)
}

func confirmClosing() {
	var conf string
	fmt.Println("Do you want to save and close the current spreadsheet? (y/s/Save - d/Discard - n/No)")
	fmt.Scan(&conf)
	conf = strings.ToLower(conf)

	switch conf {
	case "y", "yes", "s", "save":
		err := crntSheet.SaveToFile()
		if err != nil {
			fmt.Println("Error While Saving File:", err.Error())
			os.Exit(2)
		}
		close([]string{})

	case "d", "discard":
		close([]string{})

	default:
		return // n/No => No Operation
	}
}

func close(args []string) {
	checkOpenSheet()
	if len(args) > 2 {
		fmt.Println("This command takes no arguments!")
		os.Exit(1)
	}
	name := crntSheet.Name
	crntSheet = nil
	fmt.Printf("Spreadsheet (%s) has been closed.\n", name)
}

func getAt(args []string) {
	checkOpenSheet()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: get <row> <col>")
		os.Exit(1)
	}

	r, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("The row must be int!")
		os.Exit(1)
	}

	c, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("The col must be int!")
		os.Exit(1)
	}

	val, err := crntSheet.GetValueAt(r, c)
	if err != nil {
		fmt.Println("Error While Getting Value:", err.Error())
		os.Exit(1) // this is not a code.2 error because the reason mostly is a user miss-input case
	}

	// happy scenario :)
	fmt.Printf("Cell value at (%d, %d): %d\n", r, c, val)
}

func setAt(args []string) {
	checkOpenSheet()
	if len(args) != 5 {
		fmt.Println("You should follow this pattern: set <row> <col> <value>")
		os.Exit(1)
	}

	r, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("The row must be int!")
		os.Exit(1)
	}

	c, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("The col must be int!")
		os.Exit(1)
	}

	val, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println("The value must be int!")
		os.Exit(1)
	}

	err = crntSheet.SetValueAt(r, c, val)
	if err != nil {
		fmt.Println("Error While Setting Value:", err.Error())
		os.Exit(1)
	}

	fmt.Printf("New value at (%d, %d): %d\n", r, c, val)
}

func help() {
	fmt.Println("This is the help function. To be implemented...")
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Loop Interface App is not implemented yet...")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "create":
		// create <name> <rows> <cols>
		create(os.Args)
	case "open":
		open(os.Args)
	case "close":
		close(os.Args)
	case "get":
		getAt(os.Args)
	case "set":
		setAt(os.Args)
	case "help":
		help()
	default:
		fmt.Printf("Unkown command '%s'! You can see help by using 'help' command ^_^\n", cmd)
	}
}
