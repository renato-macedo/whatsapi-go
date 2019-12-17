package waconnection

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	models "github.com/renato-macedo/whatsapi/models"
	actions "github.com/renato-macedo/whatsapi/actions"
)

// Connections store all active sessions on whatsapp
var Connections []*whatsapp.Conn

// MessageHandler is responsible for handling messages
type MessageHandler struct {
	c *whatsapp.Conn
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *MessageHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

// HandleTextMessage Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (h *MessageHandler) HandleTextMessage(message whatsapp.TextMessage) {
	actions.NotifyTextMessage(message.Info.RemoteJid, message.Text, "http://localhost:3000/go")
	fmt.Printf("%v %v %v %v\n\t%v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.ContextInfo.QuotedMessageID, message.Text)
}

//HandleImageMessage Example for media handling. Video, Audio, Document are also possible in the same way
func (h *MessageHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	// fmt.Println(message)
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
			return
		}

		if _, err = h.c.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()
			if err != nil {
				return
			}
		}
	}

	userID := h.c.Info.Wid
	// create a folder with the name of the receiver
	directory := filepath.Join("./storage", userID, "images")
	fmt.Printf("RemoteJid \t %v", message.Info.RemoteJid)
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

	filename := fmt.Sprintf("%v/%v/%v/%v/%v.%v", dir, "storage", userID, "images", message.Info.Id, strings.Split(message.Type, "/")[1])
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	log.Printf("%v\timage received, saved at:%v\n", message.Info.Timestamp, filename)

}

// NewConnection start a new connection :)
func NewConnection(sessionName string, done chan models.Result) {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	var r models.Result
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}

	//Add handler
	wac.AddHandler(&MessageHandler{wac})

	//login or restore
	if err := login(wac, sessionName); err != nil {
		//log.Fatalf("error logging in: %v", err)
		r.Success = false
		r.Message = fmt.Sprintf("error logging in: %v", err)
		done <- r
		return
	}

	//verifies phone connectivity
	pong, err := wac.AdminTest()

	if !pong || err != nil {
		log.Fatalf("error pinging in: %v\n", err)
	}

	// adiciona a nova conexão no slice
	Connections = append(Connections, wac)

	// diz para a outra goroutine que tudo deu certo
	r.Success = true
	r.Message = "ok"
	done <- r

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	//Disconnect safe
	fmt.Println("Shutting down now.")
	session, err := wac.Disconnect()

	if err != nil {
		log.Fatalf("error disconnecting: %v\n", err)
	}

	if err := writeSession(session, wac.Info.Wid); err != nil {
		//log.Fatalf("error saving session: %v", err)
		log.Fatalf("error saving session: %v", err)
	}
}

func login(wac *whatsapp.Conn, sessionName string) error {
	// load saved session
	//fmt.Printf("wac %v", wac.Info.Wid)
	session, err := readSession(sessionName)
	if err == nil {
		// restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v", err)
		}
	} else {

		// no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {

			return fmt.Errorf("error during login: %v", err)
		}
		//return err
	}

	// save session
	err = writeSession(session, wac.Info.Wid)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}

	return err
}

func readSession(sessionName string) (whatsapp.Session, error) {
	session := whatsapp.Session{}
	wd, er := os.Getwd()
	if er != nil {
		fmt.Printf("Error getting the working directory %v", sessionName)
		return session, er
	}
	filename := fmt.Sprintf("%v\\%v\\%v.%v", wd, "sessions", sessionName, "gob")
	//file, err := os.Create(filename)
	file, err := os.Open(filename)
	if err != nil {
		//fmt.Printf("Error opening session file %v", err)
		return session, err
	}

	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func writeSession(session whatsapp.Session, sessionName string) error {
	dir, er := os.Getwd()
	if er != nil {
		return er
	}
	//file, err := os.Create(dir + "\\sessions\\" + sessionName + ".gob")
	filename := fmt.Sprintf("%v\\%v\\%v.%v", dir, "sessions", sessionName, "gob")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)

	err = encoder.Encode(session)

	if err != nil {
		return err
	}

	return nil
}
