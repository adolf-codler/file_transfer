// boiler {{{
package main

import (
	"adolf-codler/file_transfer/internal/cli"
	"adolf-codler/file_transfer/internal/discovery"
	"adolf-codler/file_transfer/internal/sockets"
	"adolf-codler/file_transfer/internal/transfer"
	"context"
	"log"

	//"adolf-codler/file_transfer/internal/template"
	//"bufio"
	"fmt"
	//"log"
	//"os"
	//"strings"
)

///}}}

type FileMeta struct {
	header string
}

const (
	DEFAULT_UDP_PORT = "7373"
)


// main{{{
func main() {
	mode, path, err:=cli_utils.Parse()
	if err!=nil{
		log.Fatalf("Parse Error: %v",err)
	}

	_ = transfer.ResolveData(path)

	switch mode {
	case 's':
		ctx, cancel := context.WithCancel(context.Background())
		go func(){
			if err := discovery.Broadcast(ctx, DEFAULT_UDP_PORT); err != nil{
				log.Fatalf("Broadcasting error: %s", err)
			}
		}()
		soc.StartServer(path, cancel)
	case 'r':
		ip, err:=discovery.ListenBroadcast(DEFAULT_UDP_PORT)
		if err!=nil{
			log.Fatalf("Listening Error: %s", err)
		}
		soc.StartClient(path, ip.IP.String())
	default:
		fmt.Println("Invalid choice. Please select 1 or 2.")
	}
}
//}}}
