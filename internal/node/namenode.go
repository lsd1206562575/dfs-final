package node

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	dfs_rpc "dfs/internal/rpc"
)

type NameNode struct {
}

func (c *NameNode) Example(args *dfs_rpc.ExampleArgs, reply *dfs_rpc.ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (n *NameNode) server() {
	rpc.Register(n)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	os.Remove("dfs-namenode-socket")
	l, e := net.Listen("unix", "dfs-namenode-socket")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func MakeNameNode() *NameNode{
	n := NameNode{}
	
	n.server()
	log.Printf("NameNode Server started....")

	return &n
}
