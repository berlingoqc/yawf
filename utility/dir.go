package utility

import (
	"errors"
	"io"
	"os"
)

// IsDirectoryEmpty return nil if the directory is empty
func IsDirectoryEmpty(direc string) error {
	f, err := os.Open(direc)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Readdir(1)
	if err == io.EOF {
		return nil
	}
	return errors.New(direc + " is not empty")
}
