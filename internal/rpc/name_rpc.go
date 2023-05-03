package rpc

import (
	"fmt"
	"log"
	"net/rpc"
)

func InsertFileData(args *MetaDataArgs, reply *MetaDataReply){
	callNameNode("NameNode.InsertFileData", &args, &reply)
	log.Printf("Update Redis Success!")
}

func UpdateFileData(args *UpdateMetaArgs, reply *MetaDataReply) {
	callNameNode("NameNode.UpdateFileData", &args, &reply)
	log.Printf("Update Redis Success Again.")
}

func callNameNode(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":8090")
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