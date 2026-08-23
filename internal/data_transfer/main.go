// package transfer
package transfer 

import (
	"fmt"
	"io"
	//"io/fs"
	"net"
	"os"
)

//const CHUNK=256

func RecvFile(conn net.Conn){
	defer conn.Close()
	file,err:=os.Create("/Users/adolfcodler/Downloads/Recieve")
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer file.Close()
	_,err=io.Copy(file, conn)
	if err!=nil{
		fmt.Println(err)
		return
	}
} 

func SendFile(conn net.Conn){
	file, err:=os.Open("send/porn.mp4")
	if err!=nil{
		fmt.Println(err)
	}
	defer file.Close()
	_,err=io.Copy(conn,file)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println("sent")

}
