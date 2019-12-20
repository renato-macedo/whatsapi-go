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
	"github.com/renato-macedo/whatsapi/connections"
	"github.com/renato-macedo/whatsapi/handlers"
	models "github.com/renato-macedo/whatsapi/models"
	utils "github.com/renato-macedo/whatsapi/utils"
)

// connections.Connections store all active sessions

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
	if err := login(wac, sessionName, done); err != nil {
		//log.Fatalf("error logging in: %v", err)
		r.Success = false
		//r.Message = fmt.Sprintf("error logging in: %v", err)
		r.Message = <-QRCODE
		done <- r
		return
	}
	// r.Message = <-QRCODE
	// r.Success = true
	// done <- r
	fmt.Println(r.Message)
	//verifies phone connectivity
	pong, err := wac.AdminTest()

	if !pong || err != nil {
		log.Fatalf("error pinging in: %v\n", err)
	}

	// adiciona a nova conexão no slice
	connections.Connections = append(connections.Connections, wac)

	// diz para a outra goroutine que tudo deu certo

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
	log.Printf("removing session %v", sessionName)
	// sessionName needs to be a valid whatsapp number and end with 557192665847@c.us
	connections.Connections = utils.RemoveConnection(connections.Connections, sessionName)
	fmt.Printf("%v", connections.Connections)
}

func login(wac *whatsapp.Conn, sessionName string, done chan models.Result) error {
	// load saved session
	//fmt.Printf("wac %v", wac.Info.Wid)
	var r models.Result
	session, err := readSession(sessionName)
	if err == nil {
		// restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v", err)
		}
		r.Success = false
		r.Message = "Session restored"
	} else {

		// no saved session -> regular login
		qr := make(chan string)
		go func() {
			fmt.Println("to aqui")

			r.Success = true
			r.Message = <-qr
			done <- r
			fmt.Println("agora aqui")
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
	filename := fmt.Sprintf("%v\\%v\\%v.%v", wd, "sessions", sessionName, "@c.us.gob")
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
