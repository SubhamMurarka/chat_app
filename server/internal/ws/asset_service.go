package ws

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func saveFile(fileHeader *multipart.FileHeader, key int) error {
	fileName := fileHeader.Filename
	fmt.Println(fileName)
	//TODO to get ext. to allow particular files only

	src, err := fileHeader.Open()
	if err != nil {
		return errors.New("resend image")
	}

	defer src.Close()

	//TODO generate unique file names
	path := filepath.Join("/home/murarka/chat_app/server/uploads", fileName)
	save, err := os.Create(path)
	defer save.Close()

	_, err = io.Copy(save, src)
	if err != nil {
		return errors.New("error saving image")
	}

	return nil
}
