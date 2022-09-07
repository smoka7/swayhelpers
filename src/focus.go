package src

import (
	"log"

	"github.com/joshuarubin/go-sway"
)

type Focus struct {
	hasFloating, hasTiling bool
	fullScreened           bool
	c                      Client
	tilingNodes            []*sway.Node
	floatingNodes          []*sway.Node
	focusedIndex           int
	lastIndex              int
	dir                    string
}

func newFocus(c Client, dir string) Focus {
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
		c:             c,
		dir:           dir,
	}
}

// when there is only one type of node dont modify
func (f Focus) onlyTilingOrFloating() bool {
	if f.hasFloating && f.hasTiling {
		return false
	}

	_, err := f.c.Conn.RunCommand(f.c.ctx, "focus "+f.dir)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// when focused window is on last or first toggle the mode
func (f Focus) onTheEdge() bool {
	if (f.focusedIndex == 0 && f.dir == "left") ||
		(f.focusedIndex == f.lastIndex && f.dir == "right") {
		_, err := f.c.Conn.RunCommand(f.c.ctx, "focus mode_toggle")
		if err != nil {
			log.Fatalln(err)
		}
		return true
	}

	return false
}

// on a empty workspace switch workspace
func (f Focus) switchWS() bool {
	if f.hasFloating || f.hasTiling {
		return false
	}

	dir := "next"
	if f.dir == "left" || f.dir == "bottom" {
		dir = "prev"
	}

	_, err := f.c.Conn.RunCommand(f.c.ctx, "workspace "+dir)
	if err != nil {
		log.Fatalln(err)
	}
	return true
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

	return findFirstNode(tree, func(n *sway.Node) bool {
		return n.Type == "workspace" && n.Name == focusedwsname
	})
}

func toggleFullscreen(c Client) {
	_, err := c.Conn.RunCommand(c.ctx, "fullscreen toggle")
	if err != nil {
		log.Fatalln(err)
	}
}

func (f Focus) findIndex() Focus {
	for i, node := range f.tilingNodes {
		if !node.Focused {
			continue
		}
		f.focusedIndex = i
		f.lastIndex = len(f.tilingNodes) - 1
		f.fullScreened = *node.FullscreenMode != 0
		return f
	}

	for i, node := range f.floatingNodes {
		if !node.Focused {
			continue
		}
		f.focusedIndex = i
		f.lastIndex = len(f.floatingNodes) - 1
		f.fullScreened = *node.FullscreenMode != 0
		return f
	}
	return f
}

func (c Client) Focus(dir string) {
	// finds out if focused workspace is empty or not
	f := newFocus(c, dir)

	happend := f.switchWS()
	if happend {
		return
	}

	f = f.findIndex()

	if f.fullScreened {
		toggleFullscreen(c)
		defer toggleFullscreen(c)
	}

	happend = f.onlyTilingOrFloating()
	if happend {
		return
	}

	happend = f.onTheEdge()
	if happend {
		return
	}
	// when node is tiled and its not last or first
	_, err := c.Conn.RunCommand(c.ctx, "focus "+dir)
	if err != nil {
		log.Fatalln(err)
	}
}

func getChildNodes(given []*sway.Node) (nodes []*sway.Node) {
	for _, node := range given {
		nodes = append(nodes, getAllNodesIn(node)...)
	}
	return
}
