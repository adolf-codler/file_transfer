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


type FileType int

const (
	SINGLE_FILE FileType = iota
	MULTIPLE_FILE
	FOLDER
	FILE_FOLDER
)
//const CHUNK=256

func RecvFile(conn net.Conn, path string){
	defer conn.Close()
	if path ==""{
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

func ResolveData(paths []string) FileType {
	var files, folders int
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			continue 
		}
		if info.IsDir() {
			folders++
		} else {
			files++
		}
	}

	if files == 1 && folders == 0 {
		return SINGLE_FILE
	}
	if files > 1 && folders == 0 {
		return MULTIPLE_FILE
	}
	if files == 0 && folders >= 1 {
		return FOLDER
	}
	if files >= 1 && folders >= 1 {
		return FILE_FOLDER
	}

	return SINGLE_FILE
}
