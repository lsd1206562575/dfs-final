package rpc

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type UploadArgs struct {
	FilePath string
	IPAddr   string
	Port     int
	User     string
	Password string
}

type UploadReply struct {
}
