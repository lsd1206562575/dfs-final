package client

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"time"

	"dfs/internal/meta"
	dfs_rpc "dfs/internal/rpc"
	util "dfs/internal/util"
)

type Client struct {
	// Your definitions here.
	isDone bool
}

func (c *Client) Example(args *dfs_rpc.ExampleArgs, reply *dfs_rpc.ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// 是否完成
func (c *Client) Done() bool {
	// Your code here.
	return c.isDone
}

//开启服务
func (c *Client) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":8080")
	// os.Remove("dfs-client-socket")
	// l, e := net.Listen("unix", "dfs-client-socket")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// 创建client实例
func MakeClient(filename string, method string) *Client {
	c := Client{
		isDone: false,
	}

	// Your code here.
	log.Printf("Client start.....\n")
	c.server()

	//判断走那个方法: 上传 或者 下载
	if method == "upload" {
		// 1. 得到文件的元数据信息，往 namenode 发送请求，将文件的元数据信息存入到redis中 （这个也可以远程通过调用namenode接口来获取文件的元数据）
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Get a file: %s\n", filename)
		
		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		fileSha1 := util.FileSha1(file)

		fileMeta := meta.FileMeta{
			FileName:      fileInfo.Name(),
			Tmp_Location:  filepath.Dir(filename),
			UpdatedAt: 		 time.Now().Format("2006-01-02 15:04"),
			FileSha1:      fileSha1,
			ChunkNum:     int(math.Ceil(float64(fileInfo.Size()) / (64 * 1024 * 1024))),
			FileSize: 		fileInfo.Size(),
		}
		log.Printf("File Meta data: %s, %s, %s\n", fileMeta.FileName, fileMeta.Tmp_Location, fileMeta.UpdatedAt)
		
		// 2. 把filedata 往 redis 中写
		metaArgs := dfs_rpc.MetaDataArgs{
			FileSha1   : fileMeta.FileSha1,
			FileName   : fileMeta.FileName,
			FileSize   : fileMeta.FileSize,
			ChunkNum   : fileMeta.ChunkNum,
			UpdateTime : fileMeta.UpdatedAt,
		}
		metaReply := dfs_rpc.MetaDataReply{}
		dfs_rpc.InsertFileData(&metaArgs, &metaReply)

		// 3. 将文件进行分片，并放入队列中，并将文件写入到分配的datanode中，完成后，往 namenode发送一条请求，往redis中写入数据
		fileBlockPaths := splitFile(filename)

		//远程调用datanode的上传，将分片文件上传至sftp服务器
		for _, fileblockPath := range fileBlockPaths{
			
			args := dfs_rpc.UploadArgs{
				FileBlockPath: fileblockPath,
				IPAddr: "35.236.240.242",
				Port: 2021,
				User: "admin",
				Password: "admin",
				FileSha1: fileSha1,
				Replica: 1,
			}
			reply := dfs_rpc.UploadReply{}
			dfs_rpc.UploadFileToSftp(&args, &reply)
		}
		
		// 4. 通知 namenode client即将退出
		c.isDone = true
	}

	//这里是下载方法

	return &c
}

func splitFile(infile string) []string{
		if infile == "" {
			panic("请输入正确的文件名")
		}
		
		fileInfo, err := os.Stat(infile)
		if err != nil {
			if os.IsNotExist(err) {
				panic("File doesn't exist!")
			}
			panic(err)
		}

		var chunkSize int64 = 64 * 1024 * 1024

		num := int(math.Ceil(float64(fileInfo.Size()) / (64 * 1024 * 1024)))
		filepaths := make([]string, num)

		fi, err := os.OpenFile(infile, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return filepaths
		}
		fmt.Printf("The file will be splited into %d pieces.\n", num)
		
		b := make([]byte, chunkSize)
		var i int64 = 1
		for ; i <= int64(num); i++ {
			fi.Seek((i-1) * chunkSize, 0)
			if len(b) > int(fileInfo.Size()-(i-1) * chunkSize) {
				b = make([]byte, fileInfo.Size()-(i-1) * chunkSize)
			}
			fi.Read(b)

			//TODO: filename need to be modified.
			ofile := fmt.Sprintf("D:/code/github_code/dfs/test/chunk/%s-%d.part", fileInfo.Name(), i)
			filepaths[i-1] = ofile 
			
			fmt.Printf("Create: %s\n", ofile)
			f, err := os.OpenFile(ofile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				panic(err)
			}
			f.Write(b)
			f.Close()
		}
		fi.Close()
		fmt.Println("Split Finished!")

		return filepaths
}



