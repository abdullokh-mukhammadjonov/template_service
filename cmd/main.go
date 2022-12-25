package main

import (
	"fmt"
	"net"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/abdullokh-mukhammadjonov/template_service/config"
	pb "github.com/abdullokh-mukhammadjonov/template_service/genproto/content_service"
	"github.com/abdullokh-mukhammadjonov/template_service/pkg/logger"
	"github.com/abdullokh-mukhammadjonov/template_service/service"
	"github.com/abdullokh-mukhammadjonov/template_service/service/grpc_client"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Environment, "ek_integration_service")

	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Fatal error occured", logger.Any("err-recover", err))
			os.Exit(1)
		}
	}()

	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)
	connDB := sqlx.MustConnect("postgres", psqlString)
	grpcClients, _ := grpc_client.NewGrpcClients(&cfg)

	thirdPartyService := service.ThirdPartyService(connDB, log, grpcClients, cfg)

	s := grpc.NewServer()
	pb.RegisterHandbooksServiceServer(s, thirdPartyService)

	log.Info("Listening on port", logger.String("port", cfg.RPCPort))

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error on server!", logger.Error(err))
	}

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
		panic(err)
	}
}
