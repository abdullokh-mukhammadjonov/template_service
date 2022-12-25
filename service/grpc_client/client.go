package grpc_client

import (
	"fmt"

	"github.com/abdullokh-mukhammadjonov/template_service/config"
	cs "github.com/abdullokh-mukhammadjonov/template_service/genproto/content_service"
	"google.golang.org/grpc"
)

type ServiceManager interface {
	HandbookService() cs.HandbooksServiceClient
}

type grpcClients struct {
	handbookService cs.HandbooksServiceClient
}

func (g grpcClients) HandbookService() cs.HandbooksServiceClient {
	return g.handbookService
}

func NewGrpcClients(cfg *config.Config) (ServiceManager, error) {
	connContentService, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.TemplateServiceHost, cfg.TemplateServicePort),
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)))
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		handbookService: cs.NewHandbooksServiceClient(connContentService),
	}, err
}
