package output

import (
	"io"
	"io/ioutil"
)

func GetQrImageFsWriter(namePattern string) (io.Writer, io.Closer, string, error) {
	tempFile, err := ioutil.TempFile("", namePattern)
	if err != nil {
		return nil, nil, "", err
	}

	fileName := tempFile.Name()

	return tempFile, tempFile, fileName, nil
}
