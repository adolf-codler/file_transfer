package transfer

import (
	"archive/zip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type FileType int

const (
	SINGLE_FILE FileType = iota
	MULTIPLE_FILE
)

// FileData holds metadata about the file being transferred
type FileData struct {
	Name string   `json:"name"`
	Size int64    `json:"size"`
	Type FileType `json:"type"`
}

// SendMeta sends length-prefixed JSON metadata over the connection without buffering extra bytes
func SendMeta(conn net.Conn, meta FileData) error {
	log.Printf("Sending Meta: name=%s, size=%d bytes\n", meta.Name, meta.Size)
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("SendMeta marshal error: %w", err)
	}

	// Send 4-byte big-endian length prefix
	length := uint32(len(data))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return fmt.Errorf("SendMeta length prefix error: %w", err)
	}

	// Send metadata payload
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("SendMeta data error: %w", err)
	}
	return nil
}

// RecvMeta reads exact length-prefixed JSON metadata without consuming any subsequent file stream bytes
func RecvMeta(conn net.Conn) (FileData, error) {
	var meta FileData

	// Read exact 4-byte big-endian length prefix
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return meta, fmt.Errorf("RecvMeta length prefix error: %w", err)
	}

	// Read exact metadata payload bytes
	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return meta, fmt.Errorf("RecvMeta payload error: %w", err)
	}

	if err := json.Unmarshal(buf, &meta); err != nil {
		return meta, fmt.Errorf("RecvMeta unmarshal error: %w", err)
	}

	log.Printf("Received Meta: name=%s, size=%d bytes\n", meta.Name, meta.Size)
	return meta, nil
}

func Zip(targetZipPath string, sourcePaths []string) error {// {{{
	// Create the output zip file
	zipFile, err := os.Create(targetZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Initialize a zip writer
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, sourcePath := range sourcePaths {
		err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Generate a zip header based on file info
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// Preserve the folder structure inside the zip
			// filepath.ToSlash ensures compatibility across OS (Windows uses \, zip uses /)
			header.Name = filepath.ToSlash(path)

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate // Enable compression
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			// If it's a directory, we just needed to write the header
			if info.IsDir() {
				return nil
			}

			// It's a file, so open it and copy its contents into the zip
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			return err
		})

		if err != nil {
			return err
		}
	}
	return nil
}// }}}

func Unzip(sourceZipPath, targetDir string) error {// {{{
	reader, err := zip.OpenReader(sourceZipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Get absolute target path to prevent ZipSlip attacks
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		destPath := filepath.Join(targetDir, file.Name)
		absDest, err := filepath.Abs(destPath)
		if err != nil {
			return err
		}

		// Security: Prevent ZipSlip (extracting files outside the target dir using "../../")
		if !strings.HasPrefix(absDest, absTarget+string(os.PathSeparator)) {
			// Allow exact match if targetDir is being created
			if absDest != absTarget {
				return fmt.Errorf("illegal file path (ZipSlip vulnerability): %s", destPath)
			}
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
			continue
		}

		// Ensure parent directories exist
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}// }}}

func ResolveData(paths []string) FileType {// {{{
	files := 0
	for _, p := range paths {
		if fi, err := os.Stat(p); err != nil {
			continue
		} else if fi.IsDir() {
			return MULTIPLE_FILE
		} else {
			files++
			if files > 1 {
				return MULTIPLE_FILE
			}
		}
	}
	return SINGLE_FILE
}// }}}
