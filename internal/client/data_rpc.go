package client

import (
	"fmt"
	"log"
	"net/rpc"

	dfs_rpc "dfs/internal/rpc"
)

func UploadFileToSftp(args *dfs_rpc.UploadArgs, reply *dfs_rpc.UploadReply){
	callDataNode("DataNode.UploadFileToSftp", &args, &reply)
	log.Printf("Upload Success!")
}

func callDataNode(rpcname string, args interface{}, reply interface{}) bool {
	 c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":8080")
	//c, err := rpc.DialHTTP("unix", "dfs-datanode-socket")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

