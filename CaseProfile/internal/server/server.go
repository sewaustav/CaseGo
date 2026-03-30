package server

import (
	"net/http"

	"github.com/sewaustav/CaseGoProfile/internal/db"
	"github.com/sewaustav/CaseGoProfile/internal/server/config"
	"google.golang.org/grpc"
)

type Server struct {
	DB   *db.DataBase
	HTTP *http.Server
	GRPC *grpc.Server
}

func NewServer() (*Server, error) {
	conf := config.LoadConfig()
	if conf == nil {
		panic("Config not loaded")
	}

	database := &db.DataBase{}

	if err := database.Open(
		conf.DBName,
		conf.DBUser,
		conf.DBPassword,
		conf.DBHost,
		conf.DBPort,
	); err != nil {
		return nil, err
	}

	return &Server{}, nil
}
