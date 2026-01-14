package util

import (
	"os"
	"path"
)

func CreateFolder(relativePath string, name string) error {
	folderPath := path.Join(relativePath, name)
	_, err := os.Open(folderPath)
	if err == nil {
		return nil
	}

	err = os.Mkdir(folderPath, 755)
	return err
}

func CreateBaseDirStructure() error {
	if err := CreateFolder("", "config"); err != nil {
		return err
	}
	if err := CreateFolder("", "services"); err != nil {
		return err
	}
	if err := CreateFolder("", "data"); err != nil {
		return err
	}

	return nil
}
