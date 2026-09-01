package cli_utils

import (
	"adolf-codler/file_transfer/internal/template"
	"fmt"
	"os"
)

const(
	ARG=3
)

func Parse()(byte, []string, error){
	argv:=os.Args
	argc:=len(argv)
	if argc<ARG{
		cus_template.Usage()
		return 0, []string{""}, fmt.Errorf("not enough arguments")
	}
	command:=argv[1]
	path:=argv[2:]
	if command=="s"{
		return 's', path, nil
	} else if command=="r"{
		return 'r', path, nil
	} else{
		cus_template.Usage()
		return 0, []string{""}, fmt.Errorf("invalid command") 
	}
}
