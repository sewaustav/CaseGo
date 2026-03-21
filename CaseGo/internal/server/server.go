package server

import (
	"net/http"

	"github.com/sewaustav/CaseGoCore/config"
	"github.com/sewaustav/CaseGoCore/internal/cache"
	"github.com/sewaustav/CaseGoCore/internal/cases/repository"
	service "github.com/sewaustav/CaseGoCore/internal/cases/service/core"
	"github.com/sewaustav/CaseGoCore/internal/cases/service/llm_service"
	"github.com/sewaustav/CaseGoCore/internal/db"
	"github.com/sewaustav/CaseGoCore/pkg/middleware/rs256"
)

type Server struct {
	DB   *db.DataBase
	HTTP *http.Server
}

func New() (*Server, error) {

	database := &db.DataBase{}
	conf := config.LoadConfig()
	if err := database.Open(
		conf.DBName,
		conf.DBUser,
		conf.DBPassword,
		conf.DBHost,
		conf.DBPort,
	); err != nil {
		return nil, err
	}

	caseGoRepo := repository.NewPostgresCaseRepo(database.GetDB())
	dialogRepo := repository.NewPostgresDialogRepo(database.GetDB())
	interactionsRepo := repository.NewPostgresInteractionRepo(database.GetDB())

	redisClient, err := cache.New(conf.RedisHost, conf.RedisPassword, 0)
	if err != nil {
		return nil, err
	}

	llmService := llm_service.NewLLMService(conf.LLMURL)

	caseGoService := service.NewCaseGoCoreService(redisClient, caseGoRepo, dialogRepo, interactionsRepo, llmService)

	jwtMiddleware := rs256.New(conf.PublicKey, "auth", "all")
}
