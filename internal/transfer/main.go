package transfer

import (
	"fmt"
	"io"
	"path/filepath"
	"time"
	"net"
	"os"
)


//const CHUNK=256

func RecvFile(conn net.Conn, paths []string){
	defer conn.Close()
	var path string
	if len(paths) > 0 && paths[0] != "" {
		path = paths[0]
	} else {
		path = fmt.Sprintf("receive_%s",time.Now().Format("020106_030405"))
	}
	//file,err:=os.Create(filepath.Join("/Users/adolfcodler/Downloads/Recieve",path))
	file,err:=os.Create(filepath.Join(".",path))
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

func SendFile(conn net.Conn, paths []string){
	var path string
	if len(paths) > 0 {
		path = paths[0]
	}
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

