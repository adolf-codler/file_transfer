// package transfer
package transfer

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	//"io/fs"
	"net"
	"os"
)

//const CHUNK=256

func RecvFile(conn net.Conn, path string){
	defer conn.Close()
	if path ==""{
		path = fmt.Sprintf("receive_%s",time.Now().Format("020106_030405"))
	}
	file,err:=os.Create(filepath.Join("/Users/adolfcodler/Downloads/Recieve",path))
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

func SendFile(conn net.Conn, path string){
	file, err:=os.Open(path)
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
