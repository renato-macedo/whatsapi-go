package waconnection

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	models "github.com/renato-macedo/whatsapi/models"
	"github.com/renato-macedo/whatsapi/waconnection/handlers"
)

// Connections store all active sessions
var Connections []*whatsapp.Conn

// NewConnection start a new connection :)
func NewConnection(sessionName string, done chan models.Result) {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	var r models.Result
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}

	//Add handler
	wac.AddHandler(&handlers.MessageHandler{Connection: wac})

	//login or restore
	QRCODE := make(chan string)
	if err := login(wac, sessionName, QRCODE); err != nil {
		//log.Fatalf("error logging in: %v", err)
		r.Success = false
		//r.Message = fmt.Sprintf("error logging in: %v", err)
		r.Message = <-QRCODE
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

func login(wac *whatsapp.Conn, sessionName string, QRCODE chan string) error {
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
			//terminal := qrcodeTerminal.New()
			//terminal.Get(<-qr).Print()
			QRCODE <- <-qr // kkkk
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
