package main

import (
	"fmt"
	"net"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"gitlab.udevs.io/ekadastr/ek_integration_service/config"
	pb "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/integration_service"
	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/logger"
	"gitlab.udevs.io/ekadastr/ek_integration_service/service"
	"gitlab.udevs.io/ekadastr/ek_integration_service/service/grpc_client"
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
	pb.RegisterThirdPartyServiceServer(s, thirdPartyService)

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
