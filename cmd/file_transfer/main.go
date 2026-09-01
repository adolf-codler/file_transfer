// boiler {{{
package main

import (
	"context"
	"adolf-codler/file_transfer/internal/cli"
	"adolf-codler/file_transfer/internal/discovery"
	"adolf-codler/file_transfer/internal/sockets"
	"adolf-codler/file_transfer/internal/transfer"
	//"adolf-codler/file_transfer/internal/template"
	//"bufio"
	"fmt"
	//"log"
	//"os"
	//"strings"
)

///}}}

type packet struct {
	no int
	size int
}

const (
	DEFAULT_UDP_PORT = "7373"
)

type FileType int
const(
	SINGLE_FILE FileType = iota
	MULTIPLE_FILE
	FOLDER
	FILE_FOLDER
)

// main{{{
func main() {
	mode, path, err:=cli_utils.Parse()
	if err!=nil{
		fmt.Errorf("Parse Error: %v",err)
		return
	}

	files:=transfer.ResolveData(path)

	switch mode {
	case 's':
		ctx, cancel := context.WithCancel(context.Background())
		go discovery.Broadcast(ctx, DEFAULT_UDP_PORT)
		soc.StartServer(path, cancel)
	case 'r':
		ip:=discovery.ListenBroadcast(DEFAULT_UDP_PORT)
		soc.StartClient(path, ip.IP.String())
	default:
		fmt.Println("Invalid choice. Please select 1 or 2.")
	}
}
//}}}
