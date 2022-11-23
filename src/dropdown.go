package src

import (
	"fmt"

	"github.com/joshuarubin/go-sway"
)

const (
	app_id         = "FootDropDown"
	scratchPadName = "__i3_scratch"
	termExec       = "footclient --app-id"
)

var focusedWs *sway.Node

func (c Client) Dropdown() {
	tree, _ := c.Conn.GetTree(c.ctx)
	scratch := tree.TraverseNodes(func(n *sway.Node) bool {
		return n.Name == scratchPadName
	})
	focusedWs = c.getFocusedWs()

	inFocusedWs := focusedWs.TraverseNodes(func(n *sway.Node) bool {
		return n.AppID != nil && *n.AppID == app_id+focusedWs.Name
	})

	inScratch := scratch.TraverseNodes(func(n *sway.Node) bool {
		return n.AppID != nil && *n.AppID == app_id+focusedWs.Name
	})

	if inFocusedWs == nil && inScratch == nil {
		c.execTermFor(focusedWs.Name)
	}

	if inScratch == nil {
		c.hideDropdown()
	}

	c.showFromScratch()
}

func (c Client) execTermFor(wsName string) {
	cmd := fmt.Sprintf("exec %s '%s%s'", termExec, app_id, wsName)
	c.sendCommandAndExit(cmd)
}

func (c Client) showFromScratch() {
	cmd := fmt.Sprintf("[app_id=%s%s] scratchpad show", app_id, focusedWs.Name)
	c.sendCommandAndExit(cmd)
}

func (c Client) hideDropdown() {
	cmd := fmt.Sprintf("[app_id=%s%s] move to scratchpad", app_id, focusedWs.Name)
	c.sendCommandAndExit(cmd)
}
