package server

import (
	"net/http"

	"github.com/YoungFlores/Case_Go/Profile/internal/api"
	"github.com/YoungFlores/Case_Go/Profile/internal/db"
	categoriesHandler "github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/handlers"
	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/repo"
	categoryService "github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/service"
	profileHandler "github.com/YoungFlores/Case_Go/Profile/internal/profile/handlers"
	profileRepo "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo"
	profileService "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/handlers"
	searchRepo "github.com/YoungFlores/Case_Go/Profile/internal/search/repository"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/service"
	"github.com/YoungFlores/Case_Go/Profile/pkg/middleware/rs256"
)

type Sever struct {
	HTTP *http.Server
	DB   *db.DataBase
}

func New() (*Sever, error) {

	database := &db.DataBase{}
	config := LoadConfig()

	if err := database.Open(
		config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
	); err != nil {
		return nil, err
	}

	pr := profileRepo.NewPostgresProfileRepo(database.GetDB())
	cr := repo.NewPostgresCategoryRepo(database.GetDB())
	sr := searchRepo.NewPostgresSearchRepo(database.GetDB())

	ps := profileService.NewProfileService(pr, cr)
	cs := categoryService.NewProfessionCategoryService(cr)
	ss := service.NewSearchService(sr)

	jwtMiddleware := rs256.New(config.PublicKey, "auth", "all")

	profileHandlers := profileHandler.NewProfileHandler(ps)
	categoryHandler := categoriesHandler.NewProfessionCategoryHandler(cs)
	searchHandler := handlers.NewSearchHandler(ss)

	router := api.SetupRouter(profileHandlers, searchHandler, categoryHandler, jwtMiddleware)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &Sever{
		HTTP: srv,
		DB:   database,
	}, nil
}
