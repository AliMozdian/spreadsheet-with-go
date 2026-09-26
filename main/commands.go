package main

import (
	"fmt"
	"os"
	sheet "spreadsheetCLI/spreadsheet"
	"strconv"
	"strings"
)

const POS_SEP string = ","

func checkOpenSheetExitIfNot() {
	// checks if a sheet is there is an open sheet or not
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

func checkInitExitIfNot() {
	// checks if the .sheet is initialized or not
	// use this in every command except help
	checkList := []string{ROOTD, SAVEDS, CRNT}
	for _, val := range checkList {
		_, err := os.Stat(val)
		if err != nil {
			fmt.Printf("The %s is not initialized!\nDelete .sheet/ (if exists) and restart everything with 'init' command.\n", val)
			os.Exit(1)
		}
	}
}

func tryParseIntExitIfFailed(inpValue, title string) int {
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

func tryParsePosExitIfFailed(inpPos string) (int, int) {
	// The same as tryParseIntExitIfFailed, but for dashed-positions in input, required format: <int>-<int>
	posStr := strings.Split(inpPos, POS_SEP)
	if len(posStr) != 2 {
		fmt.Printf("Positions must follow this pattern: <int>%s<int> (don't forget the '%s')\n", POS_SEP, POS_SEP)
		os.Exit(1)
	}
	r := tryParseIntExitIfFailed(posStr[0], "The row value in position-format")
	c := tryParseIntExitIfFailed(posStr[1], "The col value in position-format")

	return r, c
}

func tryParseAreaExitIfFailed(inpPos1, inpPos2 string) *sheet.Area {
	// the same as other tryParse...ExitIfFaileds, but it opens the crntSheet
	// you can acceess crntSheet by a.sh
	r1, c1 := tryParsePosExitIfFailed(inpPos1)
	r2, c2 := tryParsePosExitIfFailed(inpPos2)
	crntSheet := tryLoadCRNTExitIfFailed()

	a, err := sheet.NewArea(crntSheet, r1, c1, r2, c2)
	if err != nil {
		fmt.Println("Error While Getting Area:", err.Error())
		os.Exit(1)
	}
	return a
}

func tryLoadCRNTExitIfFailed() *sheet.Sheet {
	// Tries to load CRNT (current focused/opened sheet), in case of error, EXITS THE PRPOGRAM! (code=2)
	// Almost the same reason as the tryParseIntExitIfFailed for its existance
	crntSheet, err := loadCRNT()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	return crntSheet
}

func tryAutoSaveExitIfFailed(crntSheet *sheet.Sheet) {
	// used for auto-saveing the crntSheet at the end of modifying commands
	err := crntSheet.SaveToFile(fileName(crntSheet.Name()))
	if err != nil {
		fmt.Println("Error While Auto Saving:", err)
		os.Exit(2)
	}
}

func confirmClosing() bool {
	var conf string
	fmt.Println("Do you want to save and close the current spreadsheet? (y/s/Save - d/Discard - n/No)")
	fmt.Scan(&conf)
	conf = strings.ToLower(conf)

	switch conf {
	case "y", "yes", "s", "save":
		crnt := tryLoadCRNTExitIfFailed()
		err := crnt.SaveToFile(fileName(crnt.Name()))
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
func confirmLosingData() bool {
	// asks the user to confirm data loss is allowed to shrink the sheet, return true/false to confirm or not
	var conf string
	fmt.Println("Do you want to contract the current spreadsheet with possible data loss? (y/n)")
	fmt.Scan(&conf)
	conf = strings.ToLower(conf)

	return conf == "y" || conf == "yes"
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

func statusCmd(args []string) {
	// initialize the .sheet directory and other neccesary assets
	checkInitExitIfNot()
	if len(args) > 2 {
		fmt.Println("This command takes no arguments!")
		os.Exit(1)
	}

	crnt, err := loadCRNT()
	if err == nil {
		row, col := crnt.Size()
		fmt.Printf("Current open Sheet: '%s' size: (%d, %d)\n", crnt.Name(), row, col)
	} else {
		// a bit risky (if the err was for Error code 2), check later
		fmt.Println("Currently, there is no open Sheet to work on.")
	}
}

func displayCmd(args []string) {
	// displays the sheet in a table-like in terminal, using tabs to visualize the grid
	// extended version of statusCmd
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) > 2 {
		fmt.Println("This command takes no arguments!")
		os.Exit(1)
	}
	crntSheet := tryLoadCRNTExitIfFailed()

	rows, cols := crntSheet.Size()
	fmt.Printf("Displaying Sheet '%s' %dx%d:\n", crntSheet.Name(), rows, cols)
	fmt.Print(crntSheet.AllGridString())
}

func createCmd(args []string) {
	checkInitExitIfNot()
	if len(args) < 5 {
		fmt.Println("You should follow this pattern: create <name> <rows> <cols> --options")
		os.Exit(1)
	}

	name := args[2]
	if name == "" || name[0] == '-' {
		fmt.Println("The name cannot be empty or starts with '-'")
	}
	rows := tryParseIntExitIfFailed(args[3], "The rows")
	cols := tryParseIntExitIfFailed(args[4], "The cols")

	sh := sheet.New(rows, cols, name)

	err := sh.SaveToFile(fileName(sh.Name()))
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

			err = openSheet(sh.Name()) // the new created sheet is loaded
			if err != nil {
				fmt.Println("Error While Opening New Created Sheet:", err.Error())
				os.Exit(2)
			}
		}
	}

	fmt.Printf("Spreadsheet (%s) successfully created with %d rows and %d cols, initiated with empty cells\n", name, rows, cols)
}

func openCmd(args []string) {
	checkInitExitIfNot()
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
			return // user said no to closing CRNT, so no error (code=0)
		}
	}

	// if crntSheet != "nil" {
	// 	fmt.Printf("There is already a spreadsheet open (%s), You should close it first!\n", crntSheet.Name())
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
	fmt.Printf("Spreadsheet (%s) is open now...\n", sh.Name())
}

func closeCmd(args []string) {
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
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
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: get <row> <col>")
		os.Exit(1)
	}
	r := tryParseIntExitIfFailed(args[2], "The row")
	c := tryParseIntExitIfFailed(args[3], "The col")
	crntSheet := tryLoadCRNTExitIfFailed()

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
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 5 {
		fmt.Println("You should follow this pattern: set <row> <col> <value>")
		os.Exit(1)
	}

	r := tryParseIntExitIfFailed(args[2], "The row")
	c := tryParseIntExitIfFailed(args[3], "The col")
	val := tryParseIntExitIfFailed(args[4], "The value")
	crntSheet := tryLoadCRNTExitIfFailed()

	err := crntSheet.SetValueAt(r, c, val)
	if err != nil {
		fmt.Println("Error While Setting Value:", err.Error())
		os.Exit(1)
	}

	tryAutoSaveExitIfFailed(crntSheet)

	fmt.Printf("New value at (%d, %d): %d\n", r, c, val)
}

func delCmd(args []string) {
	// deletes a cell (clears it actually)
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: del <row> <col>")
		os.Exit(1)
	}
	r := tryParseIntExitIfFailed(args[2], "The row")
	c := tryParseIntExitIfFailed(args[3], "The col")
	crntSheet := tryLoadCRNTExitIfFailed()

	err := crntSheet.DeleteValue(r, c)
	if err != nil {
		fmt.Println("Error While Deleting Value:", err.Error())
		os.Exit(1) // this is not a code.2 error because the reason mostly is a user miss-input case
	}

	tryAutoSaveExitIfFailed(crntSheet)

	fmt.Printf("Deleted Cell at (%d, %d)\n", r, c)
}

func resizeCmd(args []string) {
	// resizes the current sheet
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: resize <rows> <cols>")
		os.Exit(1)
	}
	r := tryParseIntExitIfFailed(args[2], "The rows")
	c := tryParseIntExitIfFailed(args[3], "The cols")
	crntSheet := tryLoadCRNTExitIfFailed()
	oldR, oldC := crntSheet.Size()

	askAgain, err := crntSheet.Resize(r, c, false)
	if err != nil {
		if !askAgain {
			fmt.Println("Error While Resizing:", err.Error())
			os.Exit(2)
		}
		if !confirmLosingData() {
			fmt.Println("Resizing Aborted!")
			return
		}
		// ignore data loss
		_, err = crntSheet.Resize(r, c, true)
		if err != nil {
			fmt.Println("Erro While Resizing:", err.Error())
		}
	}

	tryAutoSaveExitIfFailed(crntSheet)

	fmt.Printf("Successfully Resized '%s' from %dx%d to %dx%d\n", crntSheet.Name(), oldR, oldC, r, c)
}

func countCmd(args []string) {
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: count <pos1> <pos2>")
		os.Exit(1)
	}
	a := tryParseAreaExitIfFailed(args[2], args[3])
	fmt.Println("The number of Valid Cells in the given Area:", a.Count())
}

func maxCmd(args []string) {
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: max <pos1> <pos2>")
		os.Exit(1)
	}
	a := tryParseAreaExitIfFailed(args[2], args[3])
	max, err := a.Max()
	if err != nil {
		fmt.Println("Error While Calculating Max:", err.Error())
	}
	fmt.Println("The Max value in the given Area:", max)
}

func minCmd(args []string) {
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: min <pos1> <pos2>")
		os.Exit(1)
	}
	a := tryParseAreaExitIfFailed(args[2], args[3])
	min, err := a.Min()
	if err != nil {
		fmt.Println("Error While Calculating min:", err.Error())
	}
	fmt.Println("The Min value in the given Area:", min)
}

func sumCmd(args []string) {
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: sum <pos1> <pos2>")
		os.Exit(1)
	}
	a := tryParseAreaExitIfFailed(args[2], args[3])
	sum, err := a.Sum()
	if err != nil {
		fmt.Println("Error While Calculating Sum:", err.Error())
	}
	fmt.Println("The Sum of values in the given Area:", sum)
}

func averageCmd(args []string) {
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()
	if len(args) != 4 {
		fmt.Println("You should follow this pattern: average <pos1> <pos2>")
		os.Exit(1)
	}
	a := tryParseAreaExitIfFailed(args[2], args[3])
	avg, err := a.Average()
	if err != nil {
		fmt.Println("Error While Calculating Average:", err.Error())
	}
	fmt.Println("The Average of values in the given Area:", avg)
}

func clearCmd(args []string) {
	// deletes all cell in the area
	checkInitExitIfNot()
	checkOpenSheetExitIfNot()

	var allCase bool = len(args) == 3 && args[2] == "all"
	if len(args) != 4 && !allCase {
		fmt.Println("You should follow this pattern: clear <pos1> <pos2>")
		fmt.Println("Or this pattern: clear all")
		os.Exit(1)
	}

	var a *sheet.Area
	if allCase {
		crntSheet := tryLoadCRNTExitIfFailed()
		a = crntSheet.TotalArea()
	} else {
		a = tryParseAreaExitIfFailed(args[2], args[3])
	}
	a.Clear()
	tryAutoSaveExitIfFailed(a.Sheet()) // a.Sheet() is crntSheet

	fmt.Println("The area is cleared! All cells in the area have been deleted.")
}

func helpCmd(args []string) {
	if len(args) > 2 {
		fmt.Println("This command takes no arguments!")
		os.Exit(1)
	}
	fmt.Println("This is the help function. To be implemented...")
}
