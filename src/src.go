package src

import (
	"context"
	"log"
	"os"

	"github.com/joshuarubin/go-sway"
)

type Client struct {
	Conn sway.Client
	ctx  context.Context
}

func NewClient() Client {
	ctx := context.Background()
	conn, _ := sway.New(ctx)
	return Client{conn, ctx}
}

// returns all the node in parent
func getAllNodesIn(parent *sway.Node) (nodes []*sway.Node) {
	if parent.PID != nil {
		nodes = append(nodes, parent)
	}

	for _, node := range parent.Nodes {
		nodes = append(nodes, getAllNodesIn(node)...)
	}

	for _, node := range parent.FloatingNodes {
		nodes = append(nodes, getAllNodesIn(node)...)
	}
	return nodes
}

func (c Client) toggleFullscreen() {
	c.sendCommand("fullscreen toggle")
}

// sendCommand  runs the command end exits when it fails
func (c Client) sendCommand(cmd string) {
	_, err := c.Conn.RunCommand(c.ctx, cmd)
	if err != nil {
		log.Fatalln(err)
	}
}

// sendCommandAndExit runs the command end exits when successful
func (c Client) sendCommandAndExit(cmd string) {
	c.sendCommand(cmd)
	os.Exit(0)
}
