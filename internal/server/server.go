package server

import (
	"net/http"

	"github.com/sewaustav/CaseGoProfile/internal/profile/repository/db"
)

type Sever struct {
	HTTP *http.Server
	DB *db.DataBase
}

func New() (*Sever, error) {

	database := &db.DataBase{}
	config := LoadConfig()

	if err := database.Open(
		config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBHost,
	); err != nil {
		return nil, err 
	}

	return &Sever{

	}, nil 
}