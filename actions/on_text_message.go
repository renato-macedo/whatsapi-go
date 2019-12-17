package actions

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// NotifyTextMessage makes a post request on the given URL
func NotifyTextMessage(number string, text string, URL string) error {
	requestBody, err := json.Marshal(map[string]string{
		"to":      number,
		"message": text,
	})

	if err != nil {
		return err
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))
	return nil
}
