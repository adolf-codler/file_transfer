// boiler {{{
package main


import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"adolf-codler/file_transfer/internal/sockets"
	"adolf-codler/file_transfer/internal/template"
)
///}}}

// main{{{
func main() {
	reader := bufio.NewReader(os.Stdin)

	template.TP()
	

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	choice := strings.TrimSpace(input)


	switch choice {
	case "1", "server", "Server":
		soc.StartServer() // Calls function from helpers.go
	case "2", "client", "Client":
		soc.StartClient() // Calls function from helpers.go
	default:
		fmt.Println("Invalid choice. Please select 1 or 2.")
	}
}
//}}}
