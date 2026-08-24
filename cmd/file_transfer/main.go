// boiler {{{
package main

import (
	"adolf-codler/file_transfer/internal/cli"
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

// main{{{
func main() {
	mode, path, ip:=cli_utils.Parse()
	if mode==0{
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
		soc.StartServer(path)
	case 'r':
		soc.StartClient(path, ip)
	default:
		fmt.Println("Invalid choice. Please select 1 or 2.")
	}
}
//}}}
