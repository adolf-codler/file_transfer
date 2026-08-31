// boiler {{{
package main

import (
	"context"
	"adolf-codler/file_transfer/internal/cli"
	"adolf-codler/file_transfer/internal/discovery"
	"adolf-codler/file_transfer/internal/sockets"

	//"adolf-codler/file_transfer/internal/data_transfer"
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

// main{{{
func main() {
	mode,path, err:=cli_utils.Parse()
	if err!=nil{
		return
	}
	/*{{{reader := bufio.NewReader(os.Stdin)

	template.TP()
	

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	choice := strings.TrimSpace(input)
	}}}*/

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
