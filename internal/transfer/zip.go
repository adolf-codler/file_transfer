package transfer

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)


type FileType int

const (
	SINGLE_FILE FileType = iota
	MULTIPLE_FILE
)
// Zip compresses a list of files or directories into a target zip file on disk.
func Zip(targetZipPath string, sourcePaths []string) error {
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
}

// Unzip extracts a zip archive to a target directory.
func Unzip(sourceZipPath, targetDir string) error {
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
}

func ResolveData(paths []string) FileType {
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
}
