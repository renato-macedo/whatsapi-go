package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadImage from the given url
func DownloadImage(URL string) (*os.File, error) {
	response, err := http.Get(URL)
	if err != nil {
		fmt.Printf("erro %v", err)
		return nil, err
	}
	defer response.Body.Close()

	slices := strings.Split(URL, "/")
	filename := slices[len(slices)-1]
	fmt.Printf("aaaa %v", filename)

	file, err := os.Create("tmp/img/" + filename)
	if err != nil {
		fmt.Printf("erro %v", err)
		return nil, err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Printf("erro %v", err)
		return nil, err
	}

	file, err = os.Open("tmp/img/" + filename)
	if err != nil {
		fmt.Printf("erro %v", err)
		return nil, err
	}
	return file, nil
}
