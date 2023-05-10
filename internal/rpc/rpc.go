package rpc

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type UploadArgs struct {
	FileBlockPath string
	SftpIpAdr     string
	FileSha1      string
	Replica       int
}

type UploadReply struct {
}

type MetaDataArgs struct {
	FileSha1   string
	FileName   string
	FileSize   int64
	ChunkNum   int
	UpdateTime string
}

type MetaDataReply struct {
}

type UpdateMetaArgs struct {
	FileSha1  string
	FileName  string
	ChunkIdx  int
	Replica   int
	SftpIpAdr string
}
