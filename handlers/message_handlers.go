package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	actions "github.com/renato-macedo/whatsapi/actions"
	"github.com/renato-macedo/whatsapi/connections"
	"github.com/renato-macedo/whatsapi/utils"
)

// MessageHandler is responsible for handling messages
type MessageHandler struct {
	Connection *whatsapp.Conn
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *MessageHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.Connection.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
			connections.Connections = utils.RemoveConnection(connections.Connections, h.Connection.Info.Wid)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
		//waconnection.Connections
		// connections.Connections = utils.RemoveConnection(connections.Connections, h.Connection.Info.Wid)
	}
}

// HandleTextMessage Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (h *MessageHandler) HandleTextMessage(message whatsapp.TextMessage) {
	actions.NotifyTextMessage(message.Info.RemoteJid, message.Text, "http://localhost:3000/go")
	fmt.Printf("%v %v %v %v\n\t%v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.ContextInfo.QuotedMessageID, message.Text)
}

//HandleImageMessage Example for media handling. Video, Audio, Document are also possible in the same way
func (h *MessageHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
			return
		}

		if _, err = h.Connection.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()

			if err != nil {
				return
			}
		}
	}
	userID := h.Connection.Info.Wid
	messageID := message.Info.Id
	messageType := strings.Split(message.Type, "/")[1]
	dirname := "images"
	saveMedia(data, userID, messageID, messageType, dirname)
}

//HandleDocumentMessage trata as mensagens do tipo documento
func (h *MessageHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	// fmt.Println(message)
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
			return
		}

		if _, err = h.Connection.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()

			if err != nil {
				return
			}
		}
	}
	userID := h.Connection.Info.Wid
	messageID := message.Info.Id
	messageType := strings.Split(message.Type, "/")[1]
	dirname := "documents"
	saveMedia(data, userID, messageID, messageType, dirname)
}

// HandleVideoMessage trata as mensagens com video
func (h *MessageHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
			return
		}

		if _, err = h.Connection.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()

			if err != nil {
				return
			}
		}
	}
	userID := h.Connection.Info.Wid
	messageID := message.Info.Id
	messageType := strings.Split(message.Type, "/")[1]
	dirname := "videos"
	saveMedia(data, userID, messageID, messageType, dirname)
}

// HandleAudioMessage trata as mensagens de audio
func (h *MessageHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
			return
		}

		if _, err = h.Connection.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()

			if err != nil {
				return
			}
		}
	}
	userID := h.Connection.Info.Wid
	messageID := message.Info.Id
	// messageType := strings.Split(message.Type, "/")[1]
	messageType := ".ogg"
	dirname := "audios"
	saveMedia(data, userID, messageID, messageType, dirname)
}

func saveMedia(data []byte, userID string, messageID string, messageType string, dirname string) {

	// create a folder with the name of the receiver
	directory := filepath.Join("./storage", userID, dirname)
	// fmt.Printf("RemoteJid \t %v", message.Info.RemoteJid)
	fmt.Println(directory)
	errOnCreate := os.MkdirAll(directory, os.ModePerm)

	if errOnCreate != nil {
		fmt.Println("Cannot create directory")
	}

	dir, er := os.Getwd()
	// fmt.Println(dir)
	if er != nil {
		er = fmt.Errorf("error %v", er)
		fmt.Printf("%v", er)

	}

	filename := fmt.Sprintf("%v/%v/%v/%v/%v.%v", dir, "storage", userID, dirname, messageID, messageType)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	log.Printf("media received, saved at:%v\n", filename)
}
