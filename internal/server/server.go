package server

import (
	"net/http"

	"github.com/sewaustav/CaseGoProfile/internal/profile/api"
	"github.com/sewaustav/CaseGoProfile/internal/profile/repository/db"
	profilerepo "github.com/sewaustav/CaseGoProfile/internal/profile/repository/profile_repo"
	profileService "github.com/sewaustav/CaseGoProfile/internal/profile/service"
	"github.com/sewaustav/CaseGoProfile/pkg/middleware"
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

	profileRepo := profilerepo.NewPostgresProfileRepo(database.GetDB())

	profileService := profileService.NewProfileService(profileRepo)

	jwtMiddleware := middleware.NewJwtAuthmiddleware()

	profileHandlers := api.NewProfileHandler(profileService)

	api := api.SetupRouter(profileHandlers, jwtMiddleware)

	srv := &http.Server{
		Addr: "8080",
		Handler: api, 
	}
	

	return &Sever{
		HTTP: srv,
		DB: database,
	}, nil 
}