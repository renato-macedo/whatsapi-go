package waconnection

import (
	whatsapp "github.com/Rhymen/go-whatsapp"
	// "github.com/Rhymen/go-whatsapp/binary/proto"
	// "log"
	"fmt"
	"os"
)

func sendTextMessage(wac *whatsapp.Conn, number string, text string) {

	message := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		Text: text,
	}

	msgID, err := wac.Send(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		os.Exit(1)
	} else {
		fmt.Println("Message Sent -> ID : " + msgID)
	}
}

func sendImageMessage(wac *whatsapp.Conn, number string, filename string) {}
