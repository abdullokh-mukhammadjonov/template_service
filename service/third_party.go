package service

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	"gitlab.udevs.io/ekadastr/ek_integration_service/config"
	es "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/entity_service"
	is "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/integration_service"
	"gitlab.udevs.io/ekadastr/ek_integration_service/modules/ek_variables"
	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/helper"
	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/logger"
	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/util"
	"gitlab.udevs.io/ekadastr/ek_integration_service/service/grpc_client"
	"gitlab.udevs.io/ekadastr/ek_integration_service/storage"
)

//order service
type thirdPartyService struct {
	storage storage.StorageI
	logger  logger.Logger
	cfg     config.Config
	clients grpc_client.ServiceManager
	is.UnimplementedThirdPartyServiceServer
}

func ThirdPartyService(db *sqlx.DB, log logger.Logger, client grpc_client.ServiceManager, cfg config.Config) *thirdPartyService {
	return &thirdPartyService{
		storage: storage.NewStoragePostgres(db),
		logger:  log,
		cfg:     cfg,
		clients: client,
	}
}

func InformWithLog(strg storage.StorageI, type_ string, indicator string, req *is.DavreestrRequest, reqBody interface{}, respBody interface{}, err error) {
	if err != nil {
		fmt.Println("InformWithLog >>> ", err)
	}

	var errOnRecord error
	if type_ == "E" {
		errOnRecord = strg.FunctionLogs().CreateIntegrationErrorLog(&is.CreateErrorLogRequest{
			Indicator:   "DR_ERR_POINT-" + indicator,
			Target:      req.EntityNumber,
			OtherTarget: fmt.Sprintf("%d", req.OrderId),
			TargetId:    req.EntityId,
			EndpointUrl: ek_variables.DavReestrUrl,
			ResBody:     err.Error(),
			ReqBody:     util.JSONStringify(reqBody),
		})
	} else if type_ == "R" {
		errOnRecord = strg.FunctionLogs().CreateIntegrationInfoLog(&is.CreateInfoLogRequest{
			Indicator:   "DR_REC_POINT-" + indicator,
			Target:      req.EntityNumber,
			OtherTarget: fmt.Sprintf("%d", req.OrderId),
			TargetId:    req.EntityId,
			EndpointUrl: ek_variables.DavReestrUrl,
			ResBody:     util.JSONStringify(respBody),
			ReqBody:     util.JSONStringify(reqBody),
		})
	}

	if errOnRecord != nil {
		fmt.Println("-------------------------------- > could not record (third party)1 <-------------------------------", errOnRecord)
	}
}

func (grpc *thirdPartyService) MakeDavreestRequest(c context.Context, req *is.DavreestrRequest) (*empty.Empty, error) {
	result := grpc.storage.ThirdParty().MakeDavreestRequest(req)
	if !result["success"].(bool) {
		InformWithLog(grpc.storage, "E", result["point"].(string), result["rpc-request"].(*is.DavreestrRequest), result["request"], result["response"], result["error"].(error))
		return &empty.Empty{}, result["error"].(error)
	}

	InformWithLog(grpc.storage, "R", result["point"].(string), result["rpc-request"].(*is.DavreestrRequest), result["request"], result["response"], nil)

	return &empty.Empty{}, nil
}

func (grpc *thirdPartyService) MakeDigitalGovernmentRequest(c context.Context, req *is.RequestByEntityIds) (*empty.Empty, error) {

	token, requestBody, err := helper.GetGovernmentAPItoken(grpc.cfg)
	if err != nil {
		fmt.Println("2, ", err)
		grpc.storage.FunctionLogs().CreateIntegrationInfoLog(&is.CreateInfoLogRequest{
			Indicator:   "ITashkentE2",
			OtherTarget: "",
			EndpointUrl: ek_variables.DavReestrUrl,
			ResBody:     err.Error(),
			ReqBody:     util.JSONStringify(requestBody),
		})
		return nil, err
	}

	for _, id := range req.EntityNumbers {
		var (
			orderId string
		)

		entity, err := grpc.clients.EntityService().Get(context.Background(), &es.ASGetRequest{Id: id})
		if err != nil {
			grpc.storage.FunctionLogs().CreateIntegrationInfoLog(&is.CreateInfoLogRequest{
				Indicator:   "ITashkentE1",
				Target:      id,
				OtherTarget: "",
				TargetId:    id,
				EndpointUrl: ek_variables.DavReestrUrl,
				ResBody:     err.Error(),
				ReqBody:     util.JSONStringify(es.ASGetRequest{Id: id}),
			})
			return nil, err
		}

		for i := len(entity.EntityProperties) - 1; i >= 0; i-- {
			if entity.EntityProperties[i].Property.Id == "622c41af652339a2a74e1f22" {
				orderId = entity.EntityProperties[i].Value
				break
			}
		}

		request, response, err := grpc.storage.ThirdParty().MakeDigitalGovernmentRequest(orderId, entity, token, grpc.cfg)
		if err != nil {
			grpc.storage.FunctionLogs().CreateIntegrationInfoLog(&is.CreateInfoLogRequest{
				Indicator:   "ITashkentE3",
				Target:      entity.EntityNumber,
				OtherTarget: orderId,
				TargetId:    id,
				EndpointUrl: ek_variables.DavReestrUrl,
				ResBody:     err.Error(),
				ReqBody:     util.JSONStringify(request),
			})
			return nil, err
		}

		grpc.storage.FunctionLogs().CreateIntegrationInfoLog(&is.CreateInfoLogRequest{
			Indicator:   "ITashkent",
			Target:      entity.EntityNumber,
			OtherTarget: orderId,
			TargetId:    id,
			EndpointUrl: ek_variables.DavReestrUrl,
			ResBody:     util.JSONStringify(response),
			ReqBody:     util.JSONStringify(request),
		})
	}

	return &empty.Empty{}, nil
}
