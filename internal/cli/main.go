package cli_utils

import (
	"adolf-codler/file_transfer/internal/template"
	"os"
)

const(
	ARG_SEND=2
	ARG_RECV=2
)

func Parse()(byte, string, string){
	argv:=os.Args
	argc:=len(argv)
	if argc<3{
		cus_template.Usage()
		return 0,"","" 
	}
	command:=argv[1]
	if command=="s"{
		return 's', argv[2], "" 
	} else if command=="r"{
		if argc<4{
			cus_template.Usage()
		}
		return 'r', argv[2], argv[3]
	} else{
		cus_template.Usage()
		return 0, "", ""
	}
}
