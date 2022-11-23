package src

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/joshuarubin/go-sway"
	"golang.org/x/term"
)

const selectorChars = "asdfgqwertzxcvbhy"

const (
	BNRM = "\x1B[40m"
	BRED = "\x1B[41m"
	BGRN = "\x1B[42m"
	BYEL = "\x1B[43m"
	BBLU = "\x1B[44m"
	BMAG = "\x1B[45m"
	BCYN = "\x1B[46m"
	BWHT = "\x1B[47m"
	FNRM = "\x1B[0m"
	FRED = "\x1B[31m"
	FGRN = "\x1B[32m"
	FYEL = "\x1B[33m"
	FBLU = "\x1B[34m"
	FMAG = "\x1B[35m"
	FCYN = "\x1B[36m"
	FWHT = "\x1B[37m"
)

var containers []*sway.Node

func (c Client) FastFocus() {
	c.getContainters()

	showContainers()

	_, err := term.MakeRaw(0)
	if err != nil {
		log.Fatalln(err)
	}

	c.readInput()
}

// reads the input when key is valid focusses that container
// or else waits for valid key none letter keys terminates the program
func (c Client) readInput() {
	in := bufio.NewReader(os.Stdin)
	for {
		r, _, err := in.ReadRune()
		if err != nil {
			log.Fatalln("stdin:", err)
		}

		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			break
		}

		c.goToContainer(r)

	}
}

// getContainters gets the focused workspace containers
func (c Client) getContainters() {
	focusedWS := c.getFocusedWs()
	containers = getAllNodesIn(focusedWS)
}

// TODO maybe show Gui
func showContainers() {
	for index, container := range containers {
		fmt.Printf(" %s %s %s %s %s\n\n", BRED, FYEL, string(selectorChars[index]), FNRM, container.Name)
	}

	fmt.Printf("%s %s", BNRM, FNRM)
}

// goToContainer switches the focus to the node that r represents
func (c Client) goToContainer(r rune) {
	index := strings.IndexRune(selectorChars, r)
	if index < 0 || index >= len(containers) {
		return
	}

	command := fmt.Sprintf("[con_id=%d] focus", containers[index].ID)
	c.sendCommandAndExit(command)
}
