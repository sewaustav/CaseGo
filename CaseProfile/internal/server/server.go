package server

import (
	"log"
	"net"
	"net/http"

	"github.com/sewaustav/CaseGoProfile/internal/api"
	grpch "github.com/sewaustav/CaseGoProfile/internal/case_profile/handlers/grpc"
	httph "github.com/sewaustav/CaseGoProfile/internal/case_profile/handlers/http"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/repository"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/service"
	"github.com/sewaustav/CaseGoProfile/internal/db"
	"github.com/sewaustav/CaseGoProfile/internal/server/config"
	"github.com/sewaustav/CaseGoProfile/pkg/middleware/rs256"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/case_go"
	"google.golang.org/grpc"
)

type Server struct {
	DB   *db.DataBase
	HTTP *http.Server
	GRPC *grpc.Server
}

func New() (*Server, error) {
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

	repo := repository.NewPostgresCaseResultRepo(database.GetDB())

	profileService := service.NewCaseResultService(repo)

	httpJwtAuthMiddleware := rs256.New(conf.PublicKey, "auth", "all")
	grpcJwtAuthMiddleware := rs256.New(conf.PublicKey, "cases", "profile")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcJwtAuthMiddleware.GrpcInterceptor),
	)

	grpcHandler := grpch.NewCaseGRPCHandler(profileService)
	httpHandler := httph.NewHttpHandler(profileService)

	pb.RegisterCasesServer(grpcServer, grpcHandler)

	httpRoutes := api.SetupRoutes(httpHandler, httpJwtAuthMiddleware)

	srv := &http.Server{
		Addr:    ":8082",
		Handler: httpRoutes,
	}

	return &Server{
		DB:   database,
		HTTP: srv,
		GRPC: grpcServer,
	}, nil
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	go func() {
		log.Println("gRPC сервер запущен на :50051")
		if err := s.GRPC.Serve(lis); err != nil {
			log.Fatalf("gRPC error: %v", err)
		}
	}()

	log.Println("HTTP сервер запущен на :8082")
	return s.HTTP.ListenAndServe()
}
