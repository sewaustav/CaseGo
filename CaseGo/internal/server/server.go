package server

import (
	"net/http"

	"github.com/sewaustav/CaseGoCore/config"
	"github.com/sewaustav/CaseGoCore/internal/api"
	"github.com/sewaustav/CaseGoCore/internal/cache"
	"github.com/sewaustav/CaseGoCore/internal/cases/handlers/grpc"
	http_handlers "github.com/sewaustav/CaseGoCore/internal/cases/handlers/http"
	"github.com/sewaustav/CaseGoCore/internal/cases/repository"
	service "github.com/sewaustav/CaseGoCore/internal/cases/service/core"
	"github.com/sewaustav/CaseGoCore/internal/cases/service/llm_service"
	"github.com/sewaustav/CaseGoCore/internal/db"
	tk "github.com/sewaustav/CaseGoCore/internal/jwt"
	"github.com/sewaustav/CaseGoCore/pkg/middleware/rs256"
)

type Server struct {
	DB    *db.DataBase
	HTTP  *http.Server
	Redis cache.Interactor
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

	tokenService := tk.NewToken(conf.PrivateKey)

	grpsClient, err := grpc.NewCaseGoGRPC(conf.GRPCSEVER, tokenService)
	if err != nil {
		return nil, err
	}

	caseGoService := service.NewCaseGoCoreService(redisClient, caseGoRepo, dialogRepo, interactionsRepo, llmService, grpsClient)

	jwtMiddleware := rs256.New(conf.PublicKey, "auth", "all")

	httpHandler := http_handlers.NewCaseGoHttpHandler(caseGoService)

	httpRoutes := api.SetupRoutes(httpHandler, jwtMiddleware)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: httpRoutes,
	}

	return &Server{
		DB:    database,
		HTTP:  srv,
		Redis: redisClient,
	}, nil
}
