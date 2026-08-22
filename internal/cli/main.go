package	cli_utils 

import(
	"os"
	"fmt"
)

func Parse(mode *byte){
	arg:=os.Args
	if len(arg)<2 || len(arg)>=4{
		fmt.Println("Usage[]")
		return 
	}
	command:=arg[1]
	if command=="s"{
		if len(arg)!=3{
			fmt.Println("choose the file to send")
			return 
		} else{
			*mode = 's'
		}
	} else if command =="r"{
		*mode='r'

	} else {
		fmt.Println("Usage[]")
		return 
	}
}
