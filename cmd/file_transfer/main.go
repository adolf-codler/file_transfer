// boiler {{{
package main

import (
	"log"
	"github.com/alecthomas/kong"
	"adolf-codler/file_transfer/internal/cli"
)///}}}


// main{{{
func main() {
	var cli cus_cli.Cli
	ctx:= kong.Parse(&cli, kong.Name("FileTransfer"), kong.Description("Transfer files over TCP"), )
	if err:=ctx.Run(); err!=nil{
		log.Fatalln("Run error: ", err)
	}
}
//}}}
