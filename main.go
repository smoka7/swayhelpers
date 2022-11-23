package main

import (
	"fmt"
	"os"

	"github.com/smoka7/swayHelper/src"
)

func main() {
	checkCommandd()
	command := os.Args[1]
	cl := src.NewClient()
	switch command {
	case "snap":
		checkDir()
		cl.Snap(os.Args[2])
	case "focus":
		checkDir()
		cl.Focus(os.Args[2])
	case "dropdown":
		cl.Dropdown()
	case "fast":
		cl.FastFocus()
	default:
		fmt.Println("enter a valid command: snap focus dropdown fast")
	}
}

func checkCommandd() {
	if len(os.Args) < 2 {
		fmt.Println("enter a command: snap focus dropdown fast")
		os.Exit(1)
	}
}

func checkDir() {
	if len(os.Args) < 3 {
		fmt.Println("enter a direction: left,right,top or bottom")
		os.Exit(1)
	}
}
