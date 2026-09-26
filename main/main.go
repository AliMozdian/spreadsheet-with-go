package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		// later: add loop.go and separate it from cli.go
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
	case "del":
		// del <row> <col>
		delCmd(os.Args)
	case "help":
		// help
		helpCmd(os.Args)
	default:
		fmt.Printf("Unkown command '%s'! You can see help by using 'help' command ^_^\n", cmd)
	}
}
