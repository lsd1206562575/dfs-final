package node

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"

	redisConn "dfs/internal/db"
	dfs_rpc "dfs/internal/rpc"
)

type NameNode struct {
}

func (c *NameNode) Example(args *dfs_rpc.ExampleArgs, reply *dfs_rpc.ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}


func (d *NameNode) Done() bool {
	return false;
}

func (n *NameNode) server() {
	rpc.Register(n)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":8090")
	//os.Remove("dfs-namenode-socket")
	//l, e := net.Listen("unix", "dfs-namenode-socket")
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

func (n *NameNode) InsertFileData(args *dfs_rpc.MetaDataArgs, reply *dfs_rpc.MetaDataReply) error{
	//获得redis连接
	rConn := redisConn.RedisPool().Get()
	defer rConn.Close()

	//将初始信息写入redis缓存
	rConn.Do("HSET", args.FileSha1, "filename", args.FileName)
	rConn.Do("HSET", args.FileSha1, "filesize", args.FileSize)
	rConn.Do("HSET", args.FileSha1, "chunkcount", args.ChunkNum)
	rConn.Do("HSET", args.FileSha1, "updateAt", args.UpdateTime)
	log.Printf("Insert Redis file data success!")
	return nil
}

func (n *NameNode) UpdateFileData(args *dfs_rpc.UpdateMetaArgs, reply *dfs_rpc.MetaDataReply) error {
	//获得redis连接
	rConn := redisConn.RedisPool().Get()
	defer rConn.Close()

	//更新
	rConn.Do("HSET", args.FileSha1, args.FileName + "_replica_"+strconv.Itoa(args.Replica), args.SftpIpAdr)
	log.Printf("Update Redis file data success!")
	return nil
}
