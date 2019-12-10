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
	model "github.com/renato-macedo/whatsapi/model"
)

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

	// create a folder with the name of the receiver
	directory := filepath.Join(".", "userid", "images")
	fmt.Println(directory)
	errOnCreate := os.MkdirAll(directory, os.ModePerm)

	if errOnCreate != nil {
		fmt.Println("Cannot create directory")
	}

	dir, er := os.Getwd()
	// fmt.Println(dir)
	if er != nil {
		fmt.Errorf("error %v", er)
		return
	}

	filename := fmt.Sprintf("%v/%v/%v/%v.%v", dir, "userid", "images", "kkkkkkkk", strings.Split(message.Type, "/")[1])
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	log.Printf("%v %v\n\timage reveived, saved at:%v\n", message.Info.Timestamp, "toptoptop", filename)

}

// NewConnection start a new connection :)
func NewConnection(username string, done chan model.Result) {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	var r model.Result
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}

	//Add handler
	wac.AddHandler(&MessageHandler{wac})

	//login or restore
	if err := login(wac, username); err != nil {
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
		r.Success = false
		r.Message = fmt.Sprintf("error disconnecting: %v\n", err)
		done <- r
		return
	}
	if err := writeSession(session, username); err != nil {
		//log.Fatalf("error saving session: %v", err)
		r.Success = false
		r.Message = fmt.Sprintf("error logging in: %v", err)
		done <- r
		return
	}
}

func login(wac *whatsapp.Conn, sessionName string) error {
	// load saved session
	session, err := readSession(sessionName)
	if err == nil {
		// restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v", err)
		}
	} else {
		fmt.Println(err)
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
	err = writeSession(session, sessionName)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}

	return err
}

func readSession(sessionName string) (whatsapp.Session, error) {
	session := whatsapp.Session{}
	dir, er := os.Getwd()
	if er != nil {
		fmt.Printf("Error opening session file 1 %v", sessionName)
		return session, er
	}
	filename := fmt.Sprintf("%v\\%v\\%v.%v", dir, "sessions", sessionName, "gob")
	fmt.Println(filename)
	///file, err := os.Create(filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening session file 2 %v", err)
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
	fmt.Println("to aqui " + dir)
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