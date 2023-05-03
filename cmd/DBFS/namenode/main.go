package main

import (
	node "dfs/internal/node"
	"time"
)

func main() {
	n := node.MakeNameNode()

	for n.Done() == false {
		time.Sleep(time.Second)
	}
}
