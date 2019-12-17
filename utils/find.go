package utils

import (
	"os"

	whatsapp "github.com/Rhymen/go-whatsapp"
)

// SessionExists verifica pelo nome se o arquivo de sessão existe na pasta
func SessionExists(sessionName string) (bool, error) {
	// // obtendo a lista de nomes do arquivos que existem na pasta sessions
	folder, err := os.Open("./sessions")
	if err != nil {
		return false, err
	}

	files, err := folder.Readdirnames(-1)
	if err != nil {
		return false, err
	}
	for _, filename := range files {
		if sessionName+"@c.us.gob" == filename {
			return true, nil
		}
	}
	return false, nil

}

// ConnectionIsActive verifica se a conexao esta ativa
func ConnectionIsActive(Connections []*whatsapp.Conn, connectionID string) (bool, error) {

	// verificando pelo nome se a sessao já existe
	for _, connection := range Connections {
		if connectionID+"@c.us" == connection.Info.Wid {
			return true, nil
		}
	}
	return false, nil
}

// FindConnectionById passa o slice de connections e o id que é o numero de telefone sem @c.us e retorna a struct da conexao
func FindConnectionByID(Connections []*whatsapp.Conn, connectionID string) *whatsapp.Conn {
	for _, connection := range Connections {
		if connectionID+"@c.us" == connection.Info.Wid {
			return connection
		}
	}
	return nil
}
