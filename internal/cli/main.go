package cli_utils

import (
	"adolf-codler/file_transfer/internal/template"
	"fmt"
	"os"
)

const(
	ARG=2
)

func Parse()(byte, []string, error){
	argv:=os.Args
	argc:=len(argv)
	if argc<ARG{
		cus_template.Usage()
		return 0, []string{""}, fmt.Errorf("not enough arguments")
	}
	command:=argv[1]
	var path []string
	if argc<3{
		path=[]string{""}
	} else{
		path=argv[2:]
	}
	if command=="s"{
		if path[0] ==""{
			return 0, path, fmt.Errorf("No file to send")
		}
		return 's', path, nil
	} else if command=="r"{
		return 'r', path, nil
	} else{
		cus_template.Usage()
		return 0, path, fmt.Errorf("invalid command") 
	}
}
