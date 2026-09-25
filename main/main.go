package main

import (
	"errors"
	"fmt"
	"os"
	sheet "spreadsheetCLI/spreadsheet"
	"strconv"
	"strings"
)

// global current sheet: the sheet that is open (loaded to Program Memory)
// var crntSheet *sheet.Sheet = nil

const ROOTD string = ".sheet/"            // the root special storage for the package (and even the app)
const CRNT string = ROOTD + "current.sht" // equivalent of crntSheet
const SAVEDS string = ROOTD + "saveds/"   // where the actual sheet files are stored (.FORMAT)
const FORMAT string = ".csv"
const PERM os.FileMode = 0700

// idea for later: add a CLI package with a Commands Struct and options (+ how many arg each one accepts)
// in that way we can use a clean general CLI manager to contribute with Sheet
// and main package will be the one using both of them by connecting the funcs in call chains

func crntSheetName() (string, error) {
	content, err := os.ReadFile(CRNT)
	if err != nil {
		return "", errors.New("Error While Reading CRNT: " + err.Error())
	}
	return string(content), nil // I hope this doesn't lead to panic :)
}

func fileName(name string) string {
	// return filename by name, hence the N for camml-case :)
	return SAVEDS + name + FORMAT
}

func checkOpenSheet() {
	// use this in sheet-required commands
	// if crntSheet == nil {
	// 	fmt.Println("There is no open spreadsheet!")
	// 	os.Exit(1)
	// }
	name, err := crntSheetName()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if name == "" {
		fmt.Println("There is no open spreadsheet!")
		os.Exit(1)
	}
}

func openSheet(name string) error {
	// actually opening sheet (writeouts name in CRNT)
	file, err := os.OpenFile(CRNT, os.O_RDWR|os.O_CREATE, PERM)
	if err != nil {
		return errors.New("Error While openning CRNT: " + err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(name)
	if err != nil {
		return errors.New("Error While Writing on CRNT: (name=" + name + ")" + err.Error())
	}
	// check later: should I check the number of bytes written on the file?
	return nil
}

func closeSheet() error {
	// actually closing sheet (clearing CRNT)
	file, err := os.OpenFile(CRNT, os.O_RDWR|os.O_CREATE, PERM)
	if err != nil {
		return errors.New("Error While Opening CRNT: " + err.Error())
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return errors.New("Error While Truncating CRNT: " + err.Error())
	}
	return nil
}

func loadCRNT() (*sheet.Sheet, error) {
	// loads the CRNT into the memory for data usage (get, set, ...)
	name, err := crntSheetName()
	if err != nil {
		return nil, errors.New("Error While Loading CRNT: " + err.Error())
	}

	crnt, err := sheet.LoadFromFile(fileName(name), name)
	if err != nil {
		return nil, errors.New("Error While Loading CRNT File: " + err.Error())
	}
	return crnt, nil
}

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

	err = sh.SaveToFile(fileName(sh.Name))
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

func confirmClosing() bool {
	var conf string
	fmt.Println("Do you want to save and close the current spreadsheet? (y/s/Save - d/Discard - n/No)")
	fmt.Scan(&conf)
	conf = strings.ToLower(conf)

	switch conf {
	case "y", "yes", "s", "save":
		crnt, err := loadCRNT()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		err = crnt.SaveToFile(fileName(crnt.Name))
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

	crntSheet, err := loadCRNT()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	val, err := crntSheet.GetValueAt(r, c)
	if err != nil {
		fmt.Println("Error While Getting Value:", err.Error())
		os.Exit(1) // this is not a code.2 error because the reason mostly is a user miss-input case
	}

	// happy scenario :)
	fmt.Printf("Cell value at (%d, %d): %d\n", r, c, val)
}

func setCmd(args []string) {
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

	crntSheet, err := loadCRNT()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	err = crntSheet.SetValueAt(r, c, val)
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

func helpCmd(args []string) {
	checkOpenSheet()
	if len(args) > 2 {
		fmt.Println("This command takes no arguments!")
		os.Exit(1)
	}
	fmt.Println("This is the help function. To be implemented...")
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Loop Interface App is not implemented yet...")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "init":
		// init
		initCmd(os.Args)
	case "create":
		// create <name> <rows> <cols>
		createCmd(os.Args)
	case "open":
		// open <name>
		openCmd(os.Args)
	case "close":
		// close
		closeCmd(os.Args)
	case "get":
		// get <row> <col>
		getCmd(os.Args)
	case "set":
		// set <row> <col> <value>
		setCmd(os.Args)
	case "help":
		// help
		helpCmd(os.Args)
	default:
		fmt.Printf("Unkown command '%s'! You can see help by using 'help' command ^_^\n", cmd)
	}
}
