package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sewaustav/Payment/internal/api"
	"github.com/sewaustav/Payment/internal/db"
	grpc_hadler "github.com/sewaustav/Payment/internal/payment/handler/grpc"
	http_handler "github.com/sewaustav/Payment/internal/payment/handler/http"
	"github.com/sewaustav/Payment/internal/payment/repository"
	service "github.com/sewaustav/Payment/internal/payment/service/api"
	"github.com/sewaustav/Payment/internal/server/config"
	"github.com/sewaustav/Payment/pkg/middleware/rs256"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	repo := repository.NewPostgresPaymentRepo(database.GetDB())
	paymentService := service.NewPaymentService(repo)
	grpcHadler := grpc_hadler.NewPaymentGRPCHandler(paymentService)
	authJwtMiddleware := rs256.New(conf.PublicKey, "auth", "all")
	httpHandler := http_handler.NewHttpHandler(paymentService)

	httpRoutes := api.SetupRoutes(httpHandler, *authJwtMiddleware)
	srv := &http.Server{
		Addr: ":8085",
		Handler: httpRoutes,
	}

	serverCert, err := tls.LoadX509KeyPair("certs/payment.crt", "certs/payment.key")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить сертификаты сервера: %w", err)
	}

	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile("certs/ca.crt")
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
	_ = grpcHadler

	return &Server{
		DB:   database,
		GRPC: grpcServer,
		HTTP: srv,
	}, nil
}