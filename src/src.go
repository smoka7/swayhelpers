package src

import (
	"context"

	"github.com/joshuarubin/go-sway"
)

type Client struct {
	Conn sway.Client
	ctx  context.Context
}

func NewClient() Client {
	ctx := context.Background()
	conn, _ := sway.New(ctx)
	return Client{
		Conn: conn,
		ctx:  ctx,
	}
}

func findFirstNode(n *sway.Node, predicate func(*sway.Node) bool) *sway.Node {
	queue := []*sway.Node{n}
	for len(queue) > 0 {
		n = queue[0]
		queue = queue[1:]

		if n == nil {
			continue
		}

		if predicate(n) {
			return n
		}

		queue = append(queue, n.Nodes...)
		queue = append(queue, n.FloatingNodes...)
	}
	return nil
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
