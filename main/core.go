package main

import (
	"errors"
	"os"
	sheet "spreadsheetCLI/spreadsheet"
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
	name := string(content) // I hope this doesn't lead to panic :)
	name = strings.Trim(name, "\n")
	return name, nil
}

func fileName(name string) string {
	// return filename by name, hence the N for camml-case :)
	return SAVEDS + name + FORMAT
}

func openSheet(name string) error {
	// actually opening sheet (writeouts name in CRNT)
	file, err := os.OpenFile(CRNT, os.O_RDWR|os.O_CREATE, PERM)
	if err != nil {
		return errors.New("Error While openning CRNT: " + err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(name + "\n")
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
