// boiler {{{
package main

import (
	//"adolf-codler/file_transfer/internal/cli"
	"adolf-codler/file_transfer/internal/sockets"
	"adolf-codler/file_transfer/internal/template"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

///}}}

// main{{{
func main() {
	//cli_utils.Parse()
	reader := bufio.NewReader(os.Stdin)

	template.TP()
	

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	choice := strings.TrimSpace(input)


	switch choice {
	case "1", "server", "Server":
		soc.StartServer() 
	case "2", "client", "Client":
		soc.StartClient()  
	default:
		fmt.Println("Invalid choice. Please select 1 or 2.")
	}
}
//}}}
