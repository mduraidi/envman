package internal

import (
	"io"
	"os"
)

// CopyTemplate copies a template file from src to dst if it doesn't exist
func CopyTemplate(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		return nil // already exists
	}
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()
	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()
	_, err = io.Copy(dstF, srcF)
	return err
}
