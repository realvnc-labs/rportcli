package output

import (
	"io"
	"os"
)

func GetQrImageFsWriter(namePattern string) (io.Writer, io.Closer, string, error) {
	tempFile, err := os.CreateTemp("", namePattern)
	if err != nil {
		return nil, nil, "", err
	}

	fileName := tempFile.Name()

	return tempFile, tempFile, fileName, nil
}
