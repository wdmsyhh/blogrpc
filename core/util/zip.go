package util

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
)

type File struct {
	Name string
	Url  string
}

func Zip(zipName string, files []File) (string, error) {
	zipFile, err := os.Create(zipName)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		resp, err := http.Get(file.Url)
		if err != nil {
			return "", err
		}
		w, err := zipWriter.Create(file.Name)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(w, resp.Body)
		resp.Body.Close()
	}

	zipWriter.Flush()
	return zipName, nil
}
