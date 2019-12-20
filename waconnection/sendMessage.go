package waconnection

import (
	whatsapp "github.com/Rhymen/go-whatsapp"
	// "github.com/Rhymen/go-whatsapp/binary/proto"
	// "log"
	"fmt"
	"os"
)

// SendTextMessage recebe a conexao, o número destinatário e o texto da mensagem
func SendTextMessage(wac *whatsapp.Conn, number string, text string) error {

	message := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		Text: text,
	}

	msgID, err := wac.Send(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		return err

		//os.Exit(1)
	}
	fmt.Println("Message Sent -> ID : " + msgID)
	return nil
}

// SendImageMessage recebe a conexao, o número destinatário e a url da imagem
func SendImageMessage(wac *whatsapp.Conn, number string, img *os.File, caption string) error {
	msg := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		Type:    "image/jpeg",
		Caption: caption,
		Content: img,
	}
	msgID, err := wac.Send(msg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		return err
	}
	fmt.Println("Message Sent -> ID : " + msgID)
	return nil

}

// SendAudioMessage etc
func SendAudioMessage(wac *whatsapp.Conn, number string, audio *os.File, isVoiceMessage bool) error {
	msg := whatsapp.AudioMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		Ptt:     isVoiceMessage,
		Type:    "audio/ogg; codecs=opus",
		Content: audio,
	}
	msgID, err := wac.Send(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		return err
	}
	fmt.Println("Message Sent -> ID : " + msgID)
	return nil
}
