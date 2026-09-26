package main

import (
	"fmt"
	"os"
	sheet "spreadsheetCLI/spreadsheet"
	"strconv"
	"strings"
)

func checkOpenSheet() {
	// use this in sheet-required commands
	name, err := crntSheetName()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	if name == "" {
		fmt.Println("There is no open spreadsheet!")
		os.Exit(1)
	}
}

func checkIntExitIfNot(inpValue, title string) int {
	// Tries parse inpValue to int, in case of rejection, EXITS THE PROGRAM! (code=1)
	// This exists because I was using this part of code very common,
	// so I siad why not have it in one line?
	val, err := strconv.Atoi(inpValue)
	if err != nil {
		fmt.Println(title, "must be of type int!")
		os.Exit(1)
	}
	return val
}

func checkLoadCRNTExitIfFailed() *sheet.Sheet {
	// Tries to load CRNT (current focused/opened sheet), in case of error, EXITS THE PRPOGRAM! (code=2)
	// Almost the same reason as the checkIntExitIfNot for its existance
	crntSheet, err := loadCRNT()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	return crntSheet
}

func confirmClosing() bool {
	var conf string
	fmt.Println("Do you want to save and close the current spreadsheet? (y/s/Save - d/Discard - n/No)")
	fmt.Scan(&conf)
	conf = strings.ToLower(conf)

	switch conf {
	case "y", "yes", "s", "save":
		crnt := checkLoadCRNTExitIfFailed()
		err := crnt.SaveToFile(fileName(crnt.Name))
		if err != nil {
			fmt.Println("Error While Saving File:", err.Error())
			os.Exit(2)
		}
		closeCmd([]string{})
		return true

	case "d", "discard":
		closeCmd([]string{})
		return true

	default:
		return false // n/No => No Operation
	}
}

//////////////////// commands used in CLI (switch in main) ////////////////////

func initCmd(args []string) {
	// initialize the .sheet directory and other neccesary assets
	if len(args) > 2 {
		fmt.Println("This command takes no arguments!")
		os.Exit(1)
	}

	_, err := os.Stat(SAVEDS)
	if err == nil {
		fmt.Println("The .sheet dir is already initialized!")
		os.Exit(1)
	}
	err = os.MkdirAll(SAVEDS, PERM)
	if err != nil {
		fmt.Println("Error While Mkdir:", err.Error())
		os.Exit(2)
	}

	file, err := os.Create(CRNT)
	if err != nil {
		fmt.Println("Error While Creating CurrentSheet:", err.Error())
		os.Exit(2)
	}
	defer file.Close()
	file.WriteString("") // for now we don't need to write anything as at init state, crnt is empty
}

func createCmd(args []string) {
	if len(args) < 5 {
		fmt.Println("You should follow this pattern: create <name> <rows> <cols> --options")
		os.Exit(1)
	}

	name := args[2]
	if name == "" || name[0] == '-' {
		fmt.Println("The name cannot be empty or starts with '-'")
	}
	rows := checkIntExitIfNot(args[3], "The rows")
	cols := checkIntExitIfNot(args[4], "The cols")

	sh := sheet.New(rows, cols, name)

	err := sh.SaveToFile(fileName(sh.Name))
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

			err = openSheet(sh.Name) // the new created sheet is loaded
			if err != nil {
				fmt.Println("Error While Opening New Created Sheet:", err.Error())
				os.Exit(2)
			}
		}
	}

	fmt.Printf("Spreadsheet (%s) successfully created with %d rows and %d cols, initiated with empty cells\n", name, rows, cols)
}

func openCmd(args []string) {
	if len(args) != 3 {
		fmt.Println("You should follow this pattern: open <name>")
		os.Exit(1)
	}

	shName, err := crntSheetName()
	if err != nil {
		fmt.Println("Error While OpenCmd:", err.Error())
		os.Exit(2)
	}

	if shName != "" {
		fmt.Printf("There is already a spreadsheet open (%s), You should close it first!\n", shName)
		if !confirmClosing() {
			fmt.Println("Aborted!")
			os.Exit(0) // user said no to closing CRNT, so no error code
		}
	}

	// if crntSheet != "nil" {
	// 	fmt.Printf("There is already a spreadsheet open (%s), You should close it first!\n", crntSheet.Name)
	// 	confirmClosing()
	// }

	name := args[2]
	// Now this is NOT neccesary! just for checking if it exists, later: with ls in saveds/
	sh, err := sheet.LoadFromFile(fileName(name), name)
	if err != nil {
		fmt.Println("Error While Loading File:", err.Error())
		os.Exit(2)
	}

	err = openSheet(name)
	if err != nil {
		fmt.Println("Error While Opening Sheet:", err.Error())
		os.Exit(2)
	}
	fmt.Printf("Spreadsheet (%s) is open now...\n", sh.Name)
}

func closeCmd(args []string) {
	checkOpenSheet()
	if len(args) > 2 {
		fmt.Println("This command takes no arguments!")
		os.Exit(1)
	}

	name, err := crntSheetName()
	if err != nil {
		fmt.Println("Error While CloseCmd:", err.Error())
	}

	err = closeSheet()
	fmt.Printf("Spreadsheet (%s) has been closed.\n", name)
}

func getCmd(args []string) {
	// printout the value of a cell
	checkOpenSheet()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: get <row> <col>")
		os.Exit(1)
	}
	r := checkIntExitIfNot(args[2], "The row")
	c := checkIntExitIfNot(args[3], "The col")
	crntSheet := checkLoadCRNTExitIfFailed()

	val, err := crntSheet.GetValueAt(r, c)
	if err != nil {
		fmt.Println("Error While Getting Value:", err.Error())
		os.Exit(1) // this is not a code.2 error because the reason mostly is a user miss-input case
	}

	// happy scenario :)
	fmt.Printf("Cell value at (%d, %d): %d\n", r, c, val)
}

func setCmd(args []string) {
	// set a cell new value (overwrite is allowed)
	checkOpenSheet()
	if len(args) != 5 {
		fmt.Println("You should follow this pattern: set <row> <col> <value>")
		os.Exit(1)
	}

	r := checkIntExitIfNot(args[2], "The row")
	c := checkIntExitIfNot(args[3], "The col")
	val := checkIntExitIfNot(args[4], "The value")
	crntSheet := checkLoadCRNTExitIfFailed()

	err := crntSheet.SetValueAt(r, c, val)
	if err != nil {
		fmt.Println("Error While Setting Value:", err.Error())
		os.Exit(1)
	}

	fmt.Printf("New value at (%d, %d): %d\n", r, c, val)

	// autosave
	err = crntSheet.SaveToFile(fileName(crntSheet.Name))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func delCmd(args []string) {
	// deletes a cell (clears it actually)
	checkOpenSheet()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: del <row> <col>")
		os.Exit(1)
	}
	r := checkIntExitIfNot(args[2], "The row")
	c := checkIntExitIfNot(args[3], "The col")
	crntSheet := checkLoadCRNTExitIfFailed()

	err := crntSheet.DeleteValue(r, c)
	if err != nil {
		fmt.Println("Error While Deleting Value:", err.Error())
		os.Exit(1) // this is not a code.2 error because the reason mostly is a user miss-input case
	}

	fmt.Printf("Deleted Cell at (%d, %d)\n", r, c)

	// autosave
	err = crntSheet.SaveToFile(fileName(crntSheet.Name))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func helpCmd(args []string) {
	if len(args) > 2 {
		fmt.Println("This command takes no arguments!")
		os.Exit(1)
	}
	fmt.Println("This is the help function. To be implemented...")
}
