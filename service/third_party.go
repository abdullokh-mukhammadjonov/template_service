package service

import (
	"github.com/jmoiron/sqlx"
	"gitlab.udevs.io/ekadastr/ek_integration_service/config"
	cs "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/content_service"
	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/logger"
	"gitlab.udevs.io/ekadastr/ek_integration_service/service/grpc_client"
	"gitlab.udevs.io/ekadastr/ek_integration_service/storage"
)

//order service
type thirdPartyService struct {
	storage storage.StorageI
	logger  logger.Logger
	cfg     config.Config
	clients grpc_client.ServiceManager
	cs.UnimplementedHandbooksServiceServer
}

func ThirdPartyService(db *sqlx.DB, log logger.Logger, client grpc_client.ServiceManager, cfg config.Config) *thirdPartyService {
	return &thirdPartyService{
		storage: storage.NewStoragePostgres(db),
		logger:  log,
		cfg:     cfg,
		clients: client,
	}
}
