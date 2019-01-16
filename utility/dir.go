package utility

import (
	"errors"
	"io"
	"os"
)

// EnsureFolderExists ensure that all folder given exists ( created if not)
func EnsureFolderExists(folders []string) error {
	for _, f := range folders {
		err := os.MkdirAll(f, 0744)
		if err != nil {
			return err
		}
	}

	return nil
}

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
