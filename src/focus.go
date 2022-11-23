package src

import (
	"github.com/joshuarubin/go-sway"
)

type Focus struct {
	hasFloating, hasTiling bool
	fullScreened           bool
	onFloating             bool
	tilingNodes            []*sway.Node
	floatingNodes          []*sway.Node
	focusedIndex           int
	lastIndex              int
	dir                    string
}

func (c Client) newFocus(dir string) Focus {
	focusedws := c.getFocusedWs()
	hasFloating := len(focusedws.FloatingNodes) > 0
	hasTiling := len(focusedws.Nodes) > 0

	// get all the nodes in the workspace
	tilingNodes := getChildNodes(focusedws.Nodes)
	floatingNodes := getChildNodes(focusedws.FloatingNodes)

	return Focus{
		hasTiling:     hasTiling,
		hasFloating:   hasFloating,
		floatingNodes: floatingNodes,
		tilingNodes:   tilingNodes,
		dir:           dir,
	}
}

// when there is only one type of node send the originall command
func (c Client) onlyTilingOrFloating(f Focus) {
	if f.hasFloating && f.hasTiling {
		return
	}
	c.sendCommandAndExit("focus " + f.dir)
}

// when focused window is on last or first toggle the mode
func (c Client) onTheEdge(f Focus) {
	// if (f.focusedIndex == 0 && f.dir == "left") ||
	// 	(f.focusedIndex == f.lastIndex && f.dir == "right") {
	// 	c.sendCommandAndExit("focus mode_toggle")
	// }
	singleFloating := f.onFloating && len(f.floatingNodes) == 1
	singleTiling := !f.onFloating && len(f.tilingNodes) == 1
	firstTiling := !f.onFloating && f.focusedIndex == 0 && f.dir != "right"
	lastTiling := !f.onFloating && f.focusedIndex == f.lastIndex && f.dir != "left"
	if singleFloating || singleTiling || firstTiling || lastTiling {
		c.sendCommandAndExit("focus mode_toggle")
	}
}

// on a empty workspace switches to next or prev workspace based on the direction
func (c Client) switchWS(f Focus) {
	if f.hasFloating || f.hasTiling {
		return
	}

	dir := "next"
	if f.dir == "left" || f.dir == "bottom" {
		dir = "prev"
	}
	c.sendCommandAndExit("workspace " + dir)
}

func (c Client) getFocusedWs() *sway.Node {
	focusedwsname := ""
	tree, _ := c.Conn.GetTree(c.ctx)
	worksapces, _ := c.Conn.GetWorkspaces(c.ctx)

	for _, workspace := range worksapces {
		if workspace.Focused {
			focusedwsname = workspace.Name
			break
		}
	}

	return tree.TraverseNodes(func(n *sway.Node) bool {
		return n.Type == sway.NodeWorkspace && n.Name == focusedwsname
	})
}

// findIndex sets the focused node index
func (f *Focus) findIndex() {
	f.lastIndex = len(f.tilingNodes) - 1
	for i, node := range f.tilingNodes {
		if !node.Focused {
			continue
		}
		f.focusedIndex = i
		f.fullScreened = node.FullscreenMode != sway.FullscreenNone
		return
	}

	f.lastIndex = len(f.floatingNodes) - 1
	for i, node := range f.floatingNodes {
		if !node.Focused {
			continue
		}
		f.onFloating = true
		f.focusedIndex = i
		f.fullScreened = node.FullscreenMode != sway.FullscreenNone
		return
	}
}

func (c Client) Focus(dir string) {
	// finds out if focused workspace is empty or not
	f := c.newFocus(dir)

	c.switchWS(f)

	f.findIndex()

	if f.fullScreened {
		c.toggleFullscreen()
		defer c.toggleFullscreen()
	}

	c.onlyTilingOrFloating(f)

	c.onTheEdge(f)

	// when node is tiled and its not last or first
	c.sendCommandAndExit("focus " + dir)
}

func getChildNodes(given []*sway.Node) (nodes []*sway.Node) {
	for _, node := range given {
		nodes = append(nodes, getAllNodesIn(node)...)
	}
	return
}
