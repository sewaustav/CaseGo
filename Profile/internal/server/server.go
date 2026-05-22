package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/YoungFlores/Case_Go/Profile/internal/api"
	"github.com/YoungFlores/Case_Go/Profile/internal/db"
	categoriesHandler "github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/handlers"
	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/repo"
	categoryService "github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/service"
	grpc_handler "github.com/YoungFlores/Case_Go/Profile/internal/profile/handlers/grpc"
	profileHandler "github.com/YoungFlores/Case_Go/Profile/internal/profile/handlers/http"
	profileRepo "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo"
	profileService "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/handlers"
	searchRepo "github.com/YoungFlores/Case_Go/Profile/internal/search/repository"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/service"
	"github.com/YoungFlores/Case_Go/Profile/pkg/middleware/rs256"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Sever struct {
	HTTP *http.Server
	DB   *db.DataBase
	GRPC *grpc.Server
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
	lr := profileRepo.NewPostgresLevelRepo(database.GetDB())

	ps := profileService.NewProfileService(pr, cr, lr)
	cs := categoryService.NewProfessionCategoryService(cr)
	ss := service.NewSearchService(sr)

	jwtMiddleware := rs256.New(config.PublicKey, "auth", "all")

	profileHandlers := profileHandler.NewProfileHandler(ps)
	categoryHandler := categoriesHandler.NewProfessionCategoryHandler(cs)
	searchHandler := handlers.NewSearchHandler(ss)

	router := api.SetupRouter(profileHandlers, searchHandler, categoryHandler, jwtMiddleware)

	grpcHandler := grpc_handler.NewProfileGrpcHandler(ps)

	serverCert, err := tls.LoadX509KeyPair("certs/profile.crt", "certs/profile.key")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить сертификаты сервера: %w", err)
	}

	certPool := x509.NewCertPool()
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ca.crt: %w", err)
	}

	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("не удалось добавить ca.crt в пул сертификатов")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13,
	}

	creds := credentials.NewTLS(tlsConfig)
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterLevelServer(grpcServer, grpcHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return &Sever{
		HTTP: srv,
		DB:   database,
		GRPC: grpcServer,
	}, nil
}
