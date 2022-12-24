package grpc_client

import (
	"fmt"

	"gitlab.udevs.io/ekadastr/ek_integration_service/config"
	dls "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/discussion_logic_service"
	es "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/entity_service"
	us "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/user_service"
	"google.golang.org/grpc"
)

type ServiceManager interface {
	EntityService() es.EntityServiceClient
	StepService() dls.StepServiceClient
	ActionHistoryService() es.ENActionHistoryServiceClient
	StaffService() us.StaffServiceClient
}

type grpcClients struct {
	entityService        es.EntityServiceClient
	stepService          dls.StepServiceClient
	actionHistoryService es.ENActionHistoryServiceClient
	staffService         us.StaffServiceClient
}

func (g grpcClients) EntityService() es.EntityServiceClient {
	return g.entityService
}

func (g grpcClients) StepService() dls.StepServiceClient {
	return g.stepService
}

func (g grpcClients) ActionHistoryService() es.ENActionHistoryServiceClient {
	return g.actionHistoryService
}
func (g grpcClients) StaffService() us.StaffServiceClient {
	return g.staffService
}
func NewGrpcClients(cfg *config.Config) (ServiceManager, error) {
	connEntityService, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.EntityServiceHost, cfg.EntityServicePort),
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)))
	if err != nil {
		return nil, err
	}

	connDiscussionLogicService, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.DiscussionLogicServiceHost, cfg.DiscussionLogicServicePort),
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)))
	if err != nil {
		return nil, err
	}

	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)))
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		entityService:        es.NewEntityServiceClient(connEntityService),
		stepService:          dls.NewStepServiceClient(connDiscussionLogicService),
		actionHistoryService: es.NewENActionHistoryServiceClient(connEntityService),
		staffService:         us.NewStaffServiceClient(connUserService),
	}, err
}
