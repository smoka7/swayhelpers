package main

import (
	"os"

	"github.com/smoka7/swayHelper/src"
)

func main() {
	cl := src.NewClient()
	command := os.Args[1]
	switch command {
	case "snap":
		arg := os.Args[2]
		cl.Snap(arg)
	case "focus":
		arg := os.Args[2]
		cl.Focus(arg)
	case "dropdown":
		cl.Dropdown()
		// case "peek":
		// 	peek(tree, arg)
	}
}
