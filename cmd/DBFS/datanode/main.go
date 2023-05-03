package main

import (
	node "dfs/internal/node"
	"time"
)

func main() {
	d := node.MakeDataNode()

	for d.Done() == false {
		time.Sleep(time.Second)
	}
}