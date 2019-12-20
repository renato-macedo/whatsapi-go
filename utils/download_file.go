package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// DownloadFile from the given url
func DownloadFile(URL string, extension string, ch chan<- *os.File) { // (*os.File, error)
	response, err := http.Get(URL)
	if err != nil {
		fmt.Printf("erro %v", err)
		// return nil, err
		ch <- nil
	}
	defer response.Body.Close()

	//slices := strings.Split(URL, "/")
	uuid := fmt.Sprintf("%v", uuid.New())
	filename := uuid + extension
	log.Printf("got file %v \n", filename)

	file, err := os.Create("tmp/media/" + filename)

	if err != nil {
		log.Printf("erro ao criar o arquivo %v \n", err)
		// return nil, err
		ch <- nil
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Printf("erro ao copiar os dados %v \n", err)
		// return nil, err
		ch <- nil
	}
	log.Println(filename)
	file, err = os.Open("tmp/media/" + filename)
	if err != nil {
		log.Printf("erro ao abrir o arquivo criado %v \n", err)
		// return nil, err
		ch <- nil
	}
	//defer file.Close()
	ch <- file
	// err = file.Close()
	// if err != nil {
	// 	log.Printf("erro ao fechar arquivo %v", err)
	// }

	//return file, nil
}
