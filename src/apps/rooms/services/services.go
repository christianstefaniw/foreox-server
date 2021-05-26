package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"server/src/constants"
	"server/src/helpers"
)

func SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	dst, err := os.Create(fmt.Sprintf("%s/%s/%s", helpers.RootDir(), constants.MEDIA_DIR, header.Filename))

	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return header.Filename, nil
}
