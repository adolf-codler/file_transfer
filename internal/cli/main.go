package	cli_utils 

import(
	"os"
	"fmt"
)

const(
	ARG_SEND=2
	ARG_RECV=2

)

func Parse()byte{
	arg:=os.Args
	if len(arg)!=2{
		fmt.Println("Usage[]")
		return 0 
	}
	command:=arg[1]
	if command=="s"{
		return 's'
	} else if command=="r"{
		return 'r'
	} else{
		return 0
	}
}
