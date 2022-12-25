package service

import (
	"github.com/abdullokh-mukhammadjonov/template_service/config"
	cs "github.com/abdullokh-mukhammadjonov/template_service/genproto/content_service"
	"github.com/abdullokh-mukhammadjonov/template_service/pkg/logger"
	"github.com/abdullokh-mukhammadjonov/template_service/service/grpc_client"
	"github.com/abdullokh-mukhammadjonov/template_service/storage"
	"github.com/jmoiron/sqlx"
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
