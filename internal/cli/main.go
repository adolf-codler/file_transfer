package cli_utils

import (
	"fmt"
	"log"
	"os"
)

const(
	ARG_SEND=2
	ARG_RECV=2

)

func Parse()(byte, string){
	argv:=os.Args
	argc:=len(argv)
	if argc<2{
		fmt.Println("Usage[]")
		return 0,"" 
	}
	command:=argv[1]
	if command=="s"{
		if argc<3{
			log.Fatalln("Filepath not found")
		}
		return 's', argv[2]
	} else if command=="r"{
		if argc<3{
			return 'r', ""
		}else{
			return 'r', argv[2]
		}
	} else{
		return 0, ""
	}
}
