package cus_template

import(
	"fmt"
)

func TP(){
	fmt.Println("==============================")
	fmt.Println("    Go Socket Application     ")
	fmt.Println("==============================")
	fmt.Println("1. Run as Server")
	fmt.Println("2. Run as Client")
	fmt.Print("Choose mode (1 or 2): ")
}
func Usage(){
	fmt.Printf("file_transfer <cmd> <file> \n cmd:s <ip>")
}
