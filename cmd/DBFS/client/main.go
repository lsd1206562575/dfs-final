package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	client "dfs/internal/client"
)

var m = flag.String("m", "", "methodName")
var f = flag.String("f", "", "fileName")

func main() {
	flag.Parse()

	if flag.NFlag() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run main.go -m methodName -f fileName\n")
		os.Exit(1)
	}

	flag.Parse()
	c := client.MakeClient(*f, *m)

	for !c.Done() {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
	log.Println("Done!")	
}