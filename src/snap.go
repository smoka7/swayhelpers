package src

import (
	"fmt"

	"github.com/joshuarubin/go-sway"
)

// sizes of container for snap
var sizes = map[string]string{
	"left":   "50ppt 100ppt",
	"right":  "50ppt 100ppt",
	"top":    "100ppt 50ppt",
	"bottom": "100ppt 50ppt",
}

// start positions of container for snap direction
var positions = map[string]string{
	"left":   "0ppt 0ppt",
	"right":  "50ppt 0ppt",
	"top":    "0ppt 0ppt",
	"bottom": "0ppt 50ppt",
}

func (c Client) Snap(dir string) {
	tree, _ := c.Conn.GetTree(c.ctx)
	focused := tree.FocusedNode()

	if focused.FullscreenMode != sway.FullscreenNone {
		c.toggleFullscreen()
	}

	cmd := fmt.Sprintf("floating enable, resize set %s , move position %s", sizes[dir], positions[dir])
	c.sendCommand(cmd)
}
