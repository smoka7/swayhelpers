package src

import (
	"log"

	"github.com/joshuarubin/go-sway"
)

func (c *Client) getFocusedWs() *sway.Node {
	focusedwsname := ""
	tree, _ := c.Conn.GetTree(c.ctx)
	w, _ := c.Conn.GetWorkspaces(c.ctx)

	for _, v := range w {
		if v.Focused {
			focusedwsname = v.Name
			break
		}
	}

	return findFirstNode(tree, func(n *sway.Node) bool {
		return n.Type == "workspace" && n.Name == focusedwsname
	})
}

func toggleFullscreen(c *Client) {
	_, err := c.Conn.RunCommand(c.ctx, "fullscreen toggle")
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Client) Focus(dir string) {
	focusedws := c.getFocusedWs()

	noFloating := len(focusedws.FloatingNodes) == 0
	noTiling := len(focusedws.Nodes) == 0

	// on a empty workspace switch workspace
	if noFloating && noTiling {
		switchWorkspace(dir, c)
		return
	}

	// get all the nodes in the workspace
	tilingNodes := getChildNodes(focusedws.Nodes)
	floatingNodes := getChildNodes(focusedws.FloatingNodes)

    fullScreened := false
	focusedIndex := -1
    lastIndex := len(floatingNodes) - 1
	findIndex := func() {
		for i, node := range tilingNodes {
			if node.Focused {
				focusedIndex = i
				lastIndex = len(tilingNodes) - 1
				fullScreened = *node.FullscreenMode != 0
				return
			}
		}
		for i, node := range floatingNodes {
			if node.Focused {
				fullScreened = *node.FullscreenMode != 0
				focusedIndex = i
			}
		}
	}

	findIndex()

	if fullScreened {
		toggleFullscreen(c)
		defer toggleFullscreen(c)
	}

	// when there is only one type of node dont modify
	if noFloating || noTiling {
		c.Conn.RunCommand(c.ctx, "focus "+dir)
		return
	}

	// when focused window is on last or first  toggle the mode
	if (focusedIndex == 0 && dir == "left") ||
		(focusedIndex == lastIndex && dir == "right") {
		c.Conn.RunCommand(c.ctx, "focus mode_toggle")
		return
	}

	// when node is tiled and its not last or first
	c.Conn.RunCommand(c.ctx, "focus "+dir)
}


func switchWorkspace(dir string, c *Client) {
	prev := dir == "left" || dir == "bottom"
	if prev {
		c.Conn.RunCommand(c.ctx, "workspace prev")
		return
	}
	c.Conn.RunCommand(c.ctx, "workspace next")
	return
}

func getChildNodes(given []*sway.Node) (nodes []*sway.Node) {
	for _, node := range given {
		nodes = append(nodes, getAllNodesIn(node)...)
	}
	return
}
