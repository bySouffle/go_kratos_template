package file

import "os"

func Exist(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.Size() == 0 || !fileInfo.IsDir() {
		return err
	}
	return nil
}
