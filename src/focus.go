package src

import (
	"fmt"
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

func (c *Client) Focus(dir string) {
	focusedIndex, FullscreenMode := -1, -1

	focusedws := c.getFocusedWs()
	windowNodes := getAllNodesIn(focusedws)
	//delete the loop and replace with find focused
	for i, v := range windowNodes {
		if v.Focused {
			focusedIndex = i
			FullscreenMode = *v.FullscreenMode
			break
		}
	}

	if FullscreenMode != 0 {
		toggleFullscreen := func() {
			_, err := c.Conn.RunCommand(c.ctx, "fullscreen toggle")
			if err != nil {
				log.Fatalln(err)
			}
		}

		toggleFullscreen()
		defer toggleFullscreen()
	}
	
	if dir == "prev" {
		_, err := c.Conn.RunCommand(c.ctx, fmt.Sprintf("[con_id=%d] focus", windowNodes[last(focusedIndex, len(windowNodes))].ID))
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	_, err := c.Conn.RunCommand(c.ctx, fmt.Sprintf("[con_id=%d] focus", windowNodes[next(focusedIndex, len(windowNodes))].ID))
	if err != nil {
		log.Fatalln(err)
	}
}

func last(index int, len int) int {
	if index == 0 {
		return len - 1
	}
	return index - 1
}

func next(index int, len int) int {
	if index == len-1 {
		return 0
	}
	return index + 1
}

func getDist(n , m sway.Node) (x, y int64) {
	getCenter:=func (n sway.Node) (x,y int64) {
	x = (n.Rect.X + n.Rect.Width)/2
	y = (n.Rect.Y + n.Rect.Height)/2
		return
	}
	cenNx,cenNy:=getCenter(n)
	cenMx,cenMy:=getCenter(m)
	x = cenNx - cenMx
	y = cenNy - cenMy
	return
}
