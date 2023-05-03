package node

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	dfs_rpc "dfs/internal/rpc"
	sftpUtil "dfs/internal/util"
)

type DataNode struct {
}

func (c *DataNode) Example(args *dfs_rpc.ExampleArgs, reply *dfs_rpc.ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (d *DataNode) Done() bool {
	return false;
}

func (d *DataNode) server() {
	rpc.Register(d)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":9000")
	// os.Remove("dfs-datanode-socket")
	//l, e := net.Listen("unix", "dfs-datanode-socket")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}


func MakeDataNode() *DataNode{
	d := DataNode{}

	d.server()
	log.Printf("DataNode server started....")
	
	return &d
}

func (d *DataNode) UploadFileToSftp(args *dfs_rpc.UploadArgs, reply *dfs_rpc.UploadReply) error{
	
	client, err := sftpUtil.Connect(args.User, args.Password, args.IPAddr, args.Port)
	
	if err != nil {
		log.Fatal("Error:", err)
		return err
	}

	//base := filepath.Base(args.FilePath)
	sftpUtil.UploadFile(client, args.FileBlockPath, "/upload")
	return nil
}
