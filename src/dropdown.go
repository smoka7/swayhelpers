package src

import (
	"fmt"
	"log"

	"github.com/joshuarubin/go-sway"
)

func (c *Client) Dropdown() {
	tree, _ := c.Conn.GetTree(c.ctx)
	scratch := findFirstNode(tree, func(n *sway.Node) bool {
		return n.Name == "__i3_scratch"
	})
	ws := c.getFocusedWs()
	dropdown := findFirstNode(tree, func(n *sway.Node) bool {
		return n.AppID != nil && *n.AppID == "FootDropDown"+ws.Name
	})
	if dropdown == nil {
		_, err := c.Conn.RunCommand(c.ctx, "exec foot --app-id \"FootDropDown"+ws.Name+"\"")
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
	inscratch := findFirstNode(scratch, func(n *sway.Node) bool {
		return n.AppID != nil && *n.AppID == "FootDropDown"+ws.Name
	})
	if inscratch == nil {
		fmt.Println("h")
		_, err := c.Conn.RunCommand(c.ctx, "[app_id=\"FootDropDown"+ws.Name+"\"] focus")
		if err != nil {
			log.Fatalln(err)
		}
		_, err = c.Conn.RunCommand(c.ctx, "move to scratchpad")

		if err != nil {
			log.Fatalln(err)
		}
		return
	}
	_, err := c.Conn.RunCommand(c.ctx, "[app_id=\"FootDropDown"+ws.Name+"\"] scratchpad show")
	if err != nil {
		log.Fatalln(err)
	}
}
