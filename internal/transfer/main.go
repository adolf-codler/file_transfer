package transfer// {{{

import (
	"adolf-codler/file_transfer/internal/progress_bar"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)// }}}

func RecvFile(conn net.Conn) error {
	defer conn.Close()

	log.Printf("Receiving Meta \r")
	meta, err := RecvMeta(conn)
	if err != nil {
		return fmt.Errorf("RecvMeta Error: %w", err)
	}

	// Always use base filename to avoid missing directory errors
	fileName := filepath.Base(meta.Name)
	if fileName == "." || fileName == "/" || fileName == "" {
		fileName = "received_file"
	}

	destPath := filepath.Join(".", fileName)
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("File creating Error: %w", err)
	}
	defer file.Close()

	pw := progress_bar.NewProgressWriter(meta.Size, "Receiving "+fileName, 500*time.Millisecond)
	target := io.MultiWriter(file, pw)

	_, err = io.CopyN(target, conn, meta.Size)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Copying Error: %w", err)
	}
	pw.Finish()
	return nil
}

func SendFile(conn net.Conn, paths []string) error {
	defer conn.Close()
	var path string
	if len(paths) > 0 {
		path = paths[0]
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("File stat Error: %w", err)
	}

	meta := FileData{
		Name: filepath.Base(path),
		Size: info.Size(),
		Type: SINGLE_FILE,
	}

	err = SendMeta(conn, meta)
	if err != nil {
		return fmt.Errorf("SendMeta Error: %w", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("File opening Error: %w", err)
	}
	defer file.Close()

	pw := progress_bar.NewProgressWriter(meta.Size, "Sending "+meta.Name, 500*time.Millisecond)
	target := io.MultiWriter(conn, pw)

	_, err = io.Copy(target, file)
	if err != nil {
		return fmt.Errorf("Copying Error: %w", err)
	}
	pw.Finish()
	return nil
}

