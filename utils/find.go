package utils

import (
	"os"
)

// SessionExists verifica pelo nome se o arquivo de sessão existe na pasta
func SessionExists(sessionName string) (bool, error) {
	// obtendo a lista de nomes do arquivos que existem na pasta sessions
	folder, err := os.Open("./sessions")
	if err != nil {
		return false, err
	}

	files, err := folder.Readdirnames(-1)
	if err != nil {
		return false, err
	}

	// verificando pelo nome se a sessao já existe
	for _, filename := range files {
		if sessionName+".gob" == filename {
			return true, nil
		}
	}
	return false, nil
}
