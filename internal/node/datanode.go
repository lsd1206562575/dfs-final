package node

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"

	dfs_rpc "dfs/internal/rpc"
	sftpUtil "dfs/internal/util"
)

type DataNode struct {
	lock sync.Locker
	cm 	 sftpUtil.SFTPConnectionManager	
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

	// Host:    "192.168.246.100",
	// Port:    2021,
	// User:    "admin",
	// Key:     "admin",

	d.cm = *sftpUtil.NewSFTPConnectionManager(10)
	log.Printf("Sftp Pool....")
	
	d.server()
	log.Printf("DataNode server started....")
	
	return &d
}

func (d *DataNode) UploadFileToSftp(args *dfs_rpc.UploadArgs, reply *dfs_rpc.UploadReply) error{

	// 使用连接池
	pool, err := d.cm.GetPool(args.SftpIpAdr, "admin", "admin")
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	client := pool.Get()
	defer pool.Put(client)
		
	// 在此处执行SFTP操作，例如client.client.ReadDir("/path/to/directory")
	sftpUtil.UploadFile(client, args.FileBlockPath, "/upload")
	return nil
}
