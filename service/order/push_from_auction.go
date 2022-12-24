package service

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/golang/protobuf/ptypes/empty"
// 	"github.com/google/uuid"
// 	"github.com/jmoiron/sqlx"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/config"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/genproto/entity_service"
// 	is "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/integration_service"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/genproto/user_service"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/modules/ek_variables"
// 	v_is "gitlab.udevs.io/ekadastr/ek_integration_service/modules/ek_variables/ek_integration_service"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/helper"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/logger"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/util"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/service/grpc_client"
// 	"gitlab.udevs.io/ekadastr/ek_integration_service/storage"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"google.golang.org/protobuf/types/known/emptypb"
// )

// type OrderService struct {
// 	storage storage.StorageI
// 	logger  logger.Logger
// 	cfg     config.Config
// 	clients grpc_client.ServiceManager
// 	is.UnimplementedOrderServiceServer
// }

// func NewOrderService(db *sqlx.DB, log logger.Logger, client grpc_client.ServiceManager, cfg config.Config) *OrderService {
// 	return &OrderService{
// 		storage: storage.NewStoragePostgres(db),
// 		logger:  log,
// 		cfg:     cfg,
// 		clients: client,
// 	}
// }

// func (grpc *OrderService) PushFromAuction(ctx context.Context, req *is.PushFromAuctionRequest) (*empty.Empty, error) {
// 	// Get this order's entity id from storage
// 	order, err := grpc.storage.Order().GetOrderByOrderId(req.OrderId)
// 	if err != nil {
// 		RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 			Indicator:  "-1",
// 			Error:      err.Error(),
// 			IsFinished: false,
// 			Success:    "NO",
// 		})
// 		return &empty.Empty{}, err
// 	}

// 	entity, err := grpc.clients.EntityService().Get(context.Background(), &entity_service.ASGetRequest{Id: order.EntityId})

// 	if err != nil {
// 		grpc.storage.ActionLog().UpdateActionLog(&v_is.UpdateActionLogRequest{ID: order.EntityId, Status: "FAILED"})
// 		RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 			EntityId:   order.EntityId,
// 			Indicator:  "-2",
// 			Error:      err.Error(),
// 			IsFinished: false,
// 			Success:    "NO",
// 		})
// 		grpc.logger.Error("Error create, update auction order", logger.Error(err))
// 		return nil, err
// 	}

// 	//                                    auksion                         yerelektron
// 	// 2: "61fa61c8a9a0964dbfcdd59a",  // Buyurtma yuborilgan          // Auksionga yuborilgan
// 	// 3: "61fbe23acd693707661d1197",  // Buyurtma qaytarilgan         // Auksiondan qaytarilgan
// 	// 4: "627e8ee75d7af0cc3032f3d9"   // Auksionga qabul qilindi      // Buyurtma qabul qilingan
// 	// 5: "61fa621da9a0964dbfcdd59c",  // Savdoga chiqarilgan          // Auksion sotuvida
// 	// 6: "61fa624fa9a0964dbfcdd59d",  // Muvaffaqiyatli yakunlandi    // Auksionda sotildi (Ro'yxatga olishda)
// 	// 7: "6234482a21d7b40ef506b888",  // Buyurtma qaytarilgan         // Auksiondan bekor qilingan
// 	statusSent := "61fa61c8a9a0964dbfcdd59a"            // 2
// 	statusRevertedAuction := "61fbe23acd693707661d1197" // 3
// 	statusAccepted := "627e8ee75d7af0cc3032f3d9"        // 4
// 	statusInAuction := "61fa621da9a0964dbfcdd59c"       // 5
// 	statusSoldInAuction := "61fa624fa9a0964dbfcdd59d"   // 6
// 	statusCanceledAuction := "6234482a21d7b40ef506b888" // 7
// 	var getProtocolResp v_is.GetAuctionProtoolResponse
// 	new_status := ""

// 	switch entity.EntityTypeCode {
// 	case 5:
// 		statusSent = "62624ad787cae4b89c431d9f"
// 		statusRevertedAuction = "62a6c60ea92f4e923db6c781"
// 		statusAccepted = "63187c21bdd563bac21131fc"
// 		statusInAuction = "62f4eecfe0506943910a48db"
// 		statusSoldInAuction = "62a6cc37a92f4e923db6cd41"
// 		statusCanceledAuction = "6305bc96b440d7b5490c072f"
// 	}

// 	statuses := map[string]string{
// 		statusSent:            "Auksionga yuborilgan",                  // 2
// 		statusRevertedAuction: "Auksiondan qaytarilgan",                // 3
// 		statusAccepted:        "Auksionga qabul qilindi",               // 4
// 		statusInAuction:       "Auksion sotuvida",                      // 5
// 		statusSoldInAuction:   "Auksionda sotildi (Ro'yxatga olishda)", // 6
// 		statusCanceledAuction: "Auksiondan bekor qilingan",             // 7
// 	}

// 	RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 		Indicator:  "-R1",
// 		Url:        req.Url,
// 		Protocol:   req.Protocol,
// 		CategoryId: req.CategoryId,
// 		OrderId:    req.OrderId,
// 	})

// 	response, reqBody, resBody, err := helper.MakeGetOrderRequest(req.OrderId, int64(entity.EntityTypeCode),
// 		grpc.cfg)

// 	if err != nil {
// 		RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 			EntityNumber: entity.EntityNumber,
// 			EntityId:     order.EntityId,
// 			PushReqRody:  reqBody,
// 			Indicator:    "-4",
// 			Error:        resBody + "  ====>  " + err.Error(),
// 			IsFinished:   false,
// 			Success:      "NO",
// 		})
// 		return &empty.Empty{}, err
// 	}

// 	if len(response.Orders) == 0 {
// 		RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 			EntityNumber: entity.EntityNumber,
// 			EntityId:     order.EntityId,
// 			Indicator:    "-5",
// 			Error:        fmt.Sprintf("no order found with this id: %d", req.OrderId),
// 			IsFinished:   false,
// 			Success:      "NO",
// 		})
// 		return &empty.Empty{}, fmt.Errorf("no order found with this id: %d", req.OrderId)
// 	}

// 	if response.Orders[0].OrderStatusesID == 3 {
// 		// if coming status is 'Auksiondan qaytarilgan', make
// 		// sure that entity's current status is not changed by mistake
// 		if entity.Status.Id == "62331cf8e936be852bfec01b" || // Xatoliklarni tahrirlashda
// 			entity.Status.Id == "62331d6be936be852bfec02c" { // Auksionga qayta yuborish
// 			RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 				EntityNumber: entity.EntityNumber,
// 				EntityId:     order.EntityId,
// 				Indicator:    "-3",
// 				Error:        "entity is not in correct status (not an error)",
// 				IsFinished:   false,
// 				Success:      "NO",
// 			})
// 			return &emptypb.Empty{}, nil
// 		}
// 	}

// 	if response.Orders[0].AcceptState == 2 {
// 		_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 			Id:       order.EntityId,
// 			StatusId: "63353f17d23d0b27354b572e", // "Yer uchastkasini qabul qilmadi" statusiga o'tadi
// 		})
// 		if err != nil {
// 			RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 				OrderId:           req.OrderId,
// 				EntityNumber:      entity.EntityNumber,
// 				EntityId:          order.EntityId,
// 				Indicator:         "-29",
// 				Error:             err.Error(),
// 				IsFinished:        false,
// 				Success:           "NO",
// 				AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 				AuctionStatusName: response.Orders[0].OrderStatus,
// 			})
// 			return &empty.Empty{}, err
// 		}
// 		_, err = grpc.clients.ActionHistoryService().Create(context.Background(), &user_service.ActionHistory{
// 			Id:             primitive.NewObjectID().Hex(),
// 			EntityId:       order.EntityId,
// 			Action:         "Auksion g'olibi yerni qabul qilmadi",
// 			EntityName:     entity.EntityNumber,
// 			StatusId:       "63353f17d23d0b27354b572e",
// 			UserId:         "000000000000000000000000",
// 			OrganizationId: "000000000000000000000000",
// 		})
// 		if err != nil {
// 			RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 				EntityNumber:      entity.EntityNumber,
// 				EntityId:          order.EntityId,
// 				Indicator:         "-30",
// 				Error:             err.Error(),
// 				IsFinished:        false,
// 				Success:           "HALF",
// 				AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 				AuctionStatusName: response.Orders[0].OrderStatus,
// 				OldStatus:         entity.Status.Name,
// 				NewStatus:         statuses[statusAccepted],
// 			})
// 			return &empty.Empty{}, err
// 		}
// 	}

// 	switch req.CategoryId {
// 	case 46, 48, 69, 72, 74, 89, 90, 92, 94, 95, 97:
// 		switch response.Orders[0].OrderStatusesID {
// 		case 2:
// 			new_status = statusSent
// 			if entity.Status.Id != statusSent {
// 				_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 					Id:       order.EntityId,
// 					StatusId: statusSent,
// 				})
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						OrderId:           req.OrderId,
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-16",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "NO",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 					})
// 					return &empty.Empty{}, err
// 				}
// 			}

// 		// order status id = 3. Auction system rejected the order with some reason
// 		// the reason is in the Description field of order [GET] request and put it to entity's property
// 		case 3:
// 			new_status = statusRevertedAuction
// 			if entity.Status.Id != statusRevertedAuction {
// 				properties := []*entity_service.EntityProperty{
// 					{PropertyId: "6242caa939c38f4fbda43748", Value: response.Orders[0].LastDescription},
// 				}
// 				_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 					Id:         order.EntityId,
// 					StatusId:   statusRevertedAuction,
// 					Properties: properties,
// 				})
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-6",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "NO",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 					})
// 					return &empty.Empty{}, err
// 				}
// 				_, err = grpc.clients.ActionHistoryService().Create(context.Background(), &user_service.ActionHistory{
// 					Id:             primitive.NewObjectID().Hex(),
// 					EntityId:       order.EntityId,
// 					Action:         "Auksiondan qaytarildi" + response.Orders[0].Description,
// 					EntityName:     "entity",
// 					UserId:         "000000000000000000000000",
// 					OrganizationId: "000000000000000000000000",
// 					StatusId:       statusRevertedAuction,
// 				})
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-7",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "HALF",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 						OldStatus:         entity.Status.Name,
// 						NewStatus:         statuses[statusRevertedAuction],
// 					})
// 					return &empty.Empty{}, err
// 				}
// 			}
// 		case 4:
// 			new_status = statusAccepted
// 			if entity.Status.Id != statusAccepted {
// 				_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 					Id:       order.EntityId,
// 					StatusId: statusAccepted,
// 				})

// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-17",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "NO",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 					})
// 					return &empty.Empty{}, err
// 				}
// 				_, err = grpc.clients.ActionHistoryService().Create(context.Background(), &user_service.ActionHistory{
// 					Id:             primitive.NewObjectID().Hex(),
// 					EntityId:       order.EntityId,
// 					Action:         "Auksionga qabul qilindi",
// 					EntityName:     "entity",
// 					StatusId:       statusAccepted,
// 					UserId:         "000000000000000000000000",
// 					OrganizationId: "000000000000000000000000",
// 				})
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-18",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "HALF",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 						OldStatus:         entity.Status.Name,
// 						NewStatus:         statuses[statusAccepted],
// 					})
// 					return &empty.Empty{}, err
// 				}
// 			}
// 		// order status id = 5. The order is assigned to auction sale on some date.
// 		// Get the auction date from order [GET] request and put it to entity's property
// 		case 5:
// 			new_status = statusInAuction
// 			if entity.Status.Id != statusInAuction {
// 				_, err = grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 					Id: order.EntityId,
// 					Properties: []*entity_service.EntityProperty{
// 						// Auksionga chiqarilgan sana
// 						{PropertyId: "628cd85818b88de7115c5abe", Value: time.Now().Format("2006-1-2")},
// 						// Buyurtma raqami
// 						{PropertyId: "61a77a0676a18e5480cc5a4a", Value: fmt.Sprint(response.Orders[0].LotNumber)},
// 						// Buksiondagi boshlang'ich narxi
// 						{PropertyId: "61c2fa298379818d937fe8f4", Value: fmt.Sprintf("%.2f", response.Orders[0].StartPrice)},
// 						// Buyurtma IDsi
// 						{PropertyId: "622c41af652339a2a74e1f22", Value: fmt.Sprint(response.Orders[0].OrderID)},
// 						// Yerning auksiondagi kategoriya idsi
// 						{PropertyId: "630da634d8e7fe958e601a1f", Value: fmt.Sprint(req.CategoryId)},
// 					},
// 					StatusId: statusInAuction,
// 				})
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-8",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "NO",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 					})
// 					return &empty.Empty{}, err
// 				}
// 				_, err = grpc.clients.ActionHistoryService().Create(context.Background(), &user_service.ActionHistory{
// 					Id:             primitive.NewObjectID().Hex(),
// 					EntityId:       order.EntityId,
// 					Action:         "Auksion savdosiga chiqarildi",
// 					EntityName:     "entity",
// 					StatusId:       statusInAuction,
// 					UserId:         "000000000000000000000000",
// 					OrganizationId: "000000000000000000000000",
// 				})
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-9",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "HALF",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 						OldStatus:         entity.Status.Name,
// 						NewStatus:         statuses[statusInAuction],
// 					})
// 					return &empty.Empty{}, err
// 				}

// 				if entity.City.Soato == 1726 {
// 					token, requestBody, err := helper.GetGovernmentAPItoken(grpc.cfg)
// 					orderIDStr := fmt.Sprintf("%d", response.Orders[0].OrderID)
// 					if err != nil {
// 						grpc.storage.FunctionLogs().CreateIntegrationInfoLog(&is.CreateInfoLogRequest{
// 							Indicator:   "ITashkentE4",
// 							Target:      entity.EntityNumber,
// 							OtherTarget: orderIDStr,
// 							TargetId:    entity.Id,
// 							EndpointUrl: ek_variables.DavReestrUrl,
// 							ResBody:     err.Error(),
// 							ReqBody:     util.JSONStringify(requestBody),
// 						})
// 					} else {
// 						grpc.storage.ThirdParty().MakeDigitalGovernmentRequest(orderIDStr, entity, token, grpc.cfg)
// 					}
// 				}
// 			}
// 		// order status id = 6. Order was sold in auction sale.
// 		// Get lot details and winner details from order [GET] request and put it to entity's property
// 		case 6:
// 			new_status = statusSoldInAuction

// 			result := response.Orders[0]
// 			lawType := "N/A"
// 			dedailsGroupNextNumber := 1
// 			sentForRegistration := false
// 			writeActionHistory := true

// 			for _, prop := range entity.EntityProperties {
// 				if prop.Property.Id == "61fa6d549cc8808579eb9924" {
// 					lawType = prop.Value
// 				} else if prop.Property.Id == "62da368c15e86da18a0f93b9" {
// 					// Auksiondan kelgan ma'lumotlar bosqichi
// 					dedailsGroupNextNumber, err = strconv.Atoi(prop.Value)
// 					if err != nil {
// 						RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 							EntityNumber:      entity.EntityNumber,
// 							EntityId:          order.EntityId,
// 							Indicator:         "-27",
// 							Error:             err.Error(),
// 							IsFinished:        false,
// 							Success:           "NO",
// 							AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 							AuctionStatusName: response.Orders[0].OrderStatus,
// 						})
// 						return &empty.Empty{}, err
// 					}
// 					dedailsGroupNextNumber++
// 				} else if prop.Property.Id == "62df8d3729523941f03bd4be" && prop.Value == "1" {
// 					// G'olibni ro'yxatga olishga yuborilgan
// 					sentForRegistration = true
// 				}
// 			}

// 			if sentForRegistration || entity.Status.Id == "61fa629ba9a0964dbfcdd59e" {
// 				return &empty.Empty{}, nil
// 			}

// 			subject_type := "Jismoniy shaxs"
// 			if result.WinnerSubjectType == 0 {
// 				subject_type = "Yuridik shaxs"
// 			}

// 			newDetails := map[string]string{
// 				"630da634d8e7fe958e601a1f": fmt.Sprint(req.CategoryId),                         // Yerning auksiondagi kategoriya idsi
// 				"61443cedfed246d6d74f0014": result.WinnerName,                                  // Auksion g'olibi
// 				"61a77a0676a18e5480cc5a4a": fmt.Sprint(response.Orders[0].LotNumber),           //buyurtma raqami
// 				"61c2fa298379818d937fe8f4": fmt.Sprintf("%.2f", response.Orders[0].StartPrice), //auksiondagi boshlang'ich narxi
// 				"61443cedfed246d6d74f001b": fmt.Sprintf("%.2f", result.SoldPrice),              // Yer uchaskasini sotilgan bahosi
// 				"622c41af652339a2a74e1f22": fmt.Sprint(result.OrderID),                         // Buyurtma IDsi
// 				"61a783d576a18e5480cc5a4b": fmt.Sprint(result.WinnerSubjectType),               // Auksion g'olibi: Subyekt turi
// 				"622b391b27ec4fc9e60f287f": result.WinnerPassport,                              // Auksion g'olibi passport seriya va raqami
// 				"622c3cbc652339a2a74e1f21": result.ProtocolFileUrl,                             // Lot fayli
// 				"622c4293dc61daec3adfbfb4": result.WinnerPinfl,                                 // Auksion g'olibi pinfl
// 				"622c4638dc61daec3adfbfe0": result.WinnerInn,                                   // Auksion g'olibing inn raqami
// 				"61443cedfed246d6d74f0016": result.AuctionDate,                                 // Auksion o'tkaziladigan sana
// 				"62ac81e99b37baa916216d8d": result.WinnerPhone,                                 // Auksion g'olibi telefon raqami
// 				"622b3a71ce7aaf40596dd160": result.WinnerAddress,                               // Auksion g'olibi manzili
// 				"622b41f6ce7aaf40596dd27a": "not.provided.by.auction.api@mail.ru",              // Auksion g'olibi elektron manzili
// 			}

// 			new_properties := []*entity_service.EntityProperty{
// 				// Auksion g'olibi
// 				{PropertyId: "61443cedfed246d6d74f0014", Value: result.WinnerName},
// 				//buyurtma raqami
// 				{PropertyId: "61a77a0676a18e5480cc5a4a", Value: fmt.Sprint(response.Orders[0].LotNumber)},
// 				//auksiondagi boshlang'ich narxi
// 				{PropertyId: "61c2fa298379818d937fe8f4", Value: fmt.Sprintf("%.2f", response.Orders[0].StartPrice)},
// 				// Yer uchaskasini sotilgan bahosi
// 				{PropertyId: "61443cedfed246d6d74f001b", Value: fmt.Sprintf("%.2f", result.SoldPrice)},
// 				// Buyurtma IDsi
// 				{PropertyId: "622c41af652339a2a74e1f22", Value: fmt.Sprint(result.OrderID)},
// 				// Auksion g'olibi: Subyekt turi
// 				{PropertyId: "61a783d576a18e5480cc5a4b", Value: fmt.Sprint(result.WinnerSubjectType), ValueLabel: subject_type},
// 				// Auksion g'olibi passport seriya va raqami
// 				{PropertyId: "622b391b27ec4fc9e60f287f", Value: result.WinnerPassport},
// 				// Lot fayli
// 				{PropertyId: "622c3cbc652339a2a74e1f21", Value: result.ProtocolFileUrl},
// 				// Auksion g'olibi pinfl
// 				{PropertyId: "622c4293dc61daec3adfbfb4", Value: result.WinnerPinfl},
// 				// Auksion g'olibing inn raqami
// 				{PropertyId: "622c4638dc61daec3adfbfe0", Value: result.WinnerInn},
// 				// Auksion o'tkaziladigan sana
// 				{PropertyId: "61443cedfed246d6d74f0016", Value: result.AuctionDate},
// 				// Auksion g'olibi telefon raqami
// 				{PropertyId: "62ac81e99b37baa916216d8d", Value: result.WinnerPhone},
// 				// Auksion g'olibi manzili
// 				{PropertyId: "622b3a71ce7aaf40596dd160", Value: result.WinnerAddress},
// 				// Auksion g'olibi elektron manzili
// 				{PropertyId: "622b41f6ce7aaf40596dd27a", Value: "not.provided.by.auction.api@mail.ru"},
// 				// Auksiondan kelgan ma'lumotlar bosqichi
// 				{PropertyId: "62da368c15e86da18a0f93b9", Value: fmt.Sprint(dedailsGroupNextNumber)},
// 			}

// 			readyForRegistration := false

// 			if lawType == "2" {
// 				temp, _, _, err := helper.MakeGetProtocolRequest(req.OrderId, int64(entity.EntityTypeCode), grpc.cfg)
// 				getProtocolResp = temp
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-22",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "NO",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 					})
// 					return &empty.Empty{}, err
// 				}

// 				if getProtocolResp.FileURL != "" {
// 					newDetails["62848d7c0effe0b32f2aff43"] = getProtocolResp.FileURL // Auksion bayonnomasi
// 					newDetails["62848d1e0effe0b32f2afed8"] = "1"                     // Auksion protokoli
// 					new_properties = append(new_properties,
// 						&entity_service.EntityProperty{
// 							// Auksion bayonnomasi
// 							PropertyId: "62848d7c0effe0b32f2aff43",
// 							Value:      getProtocolResp.FileURL,
// 						},
// 						&entity_service.EntityProperty{
// 							// Auksion protokili
// 							PropertyId: "62848d1e0effe0b32f2afed8",
// 							Value:      "1",
// 						},
// 					)
// 					if getProtocolResp.ContractFileURL != "" && getProtocolResp.ContractDate != "" {
// 						readyForRegistration = true
// 						newDetails["61443cedfed246d6d74f001f"] = getProtocolResp.ContractNumber  // Shartnoma raqami
// 						newDetails["61443cedfed246d6d74f001e"] = getProtocolResp.ContractDate    // Shartnoma sanasi
// 						newDetails["62cbc7058a472d49129c2d44"] = getProtocolResp.ContractFileURL // Ijara shartnomasi uchun havola
// 						new_properties = append(new_properties,
// 							&entity_service.EntityProperty{
// 								// Shartnoma raqami
// 								PropertyId: "61443cedfed246d6d74f001f",
// 								Value:      getProtocolResp.ContractNumber,
// 							},
// 							&entity_service.EntityProperty{
// 								// Shartnoma sanasi
// 								PropertyId: "61443cedfed246d6d74f001e",
// 								Value:      getProtocolResp.ContractDate,
// 							},
// 							&entity_service.EntityProperty{
// 								// Ijara shartnomasi uchun havola
// 								PropertyId: "62cbc7058a472d49129c2d44",
// 								Value:      getProtocolResp.ContractFileURL,
// 							})
// 					}
// 				} else {
// 					newDetails["62848d7c0effe0b32f2aff43"] = req.Url                  // Auksion bayonnomasi
// 					newDetails["62848d1e0effe0b32f2afed8"] = fmt.Sprint(req.Protocol) // Auksion protokili
// 					new_properties = append(new_properties,
// 						&entity_service.EntityProperty{
// 							// Auksion bayonnomasi
// 							PropertyId: "62848d7c0effe0b32f2aff43",
// 							Value:      req.Url,
// 						},
// 						&entity_service.EntityProperty{
// 							// Auksion protokili
// 							PropertyId: "62848d1e0effe0b32f2afed8",
// 							Value:      fmt.Sprint(req.Protocol),
// 						},
// 					)
// 				}
// 			}

// 			if helper.GotDifferentAuctionPushDetails(newDetails, entity.EntityProperties) {
// 				recordMessage := "Auksion savdosida sotildi"

// 				if entity.Status.Id == "61fa624fa9a0964dbfcdd59d" || entity.Status.Id == "62be78e1542bd77ff4b0b8b7" {
// 					// Auksionda sotildi (Ro'yxatga olishda)
// 					recordMessage = fmt.Sprintf("Auksion ma'lumotlar qabul qilindi (%d-bosqich)", dedailsGroupNextNumber)
// 				}

// 				_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 					Id:         order.EntityId,
// 					Properties: new_properties,
// 					StatusId:   statusSoldInAuction,
// 				})

// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-19",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "HALF",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 						OldStatus:         entity.Status.Name,
// 						NewStatus:         statuses[statusSoldInAuction],
// 					})
// 					return &empty.Empty{}, err
// 				}

// 				if lawType == "2" && readyForRegistration {
// 					// law_type is (2) "Ijara huquqi"
// 					recordMessage = ek_variables.DavReestrRentalRegistrationMessage

// 					// if dedailsGroupNextNumber == 1 {

// 					// }

// 					fetchAndZipReq := v_is.FetchURLsAndZip{
// 						Urls: []map[string]string{
// 							{
// 								"url": getProtocolResp.ContractFileURL,
// 								"as":  "shartnoma",
// 							},
// 							{
// 								"url": getProtocolResp.FileURL,
// 								"as":  "bayonnoma",
// 							},
// 						},
// 						OrderID: req.OrderId,
// 					}
// 					zipFileLocation, err := helper.FetchFilesAndZip(fetchAndZipReq.Urls, fmt.Sprintf("protocol_bayonnoma/%d/", fetchAndZipReq.OrderID))

// 					if err != nil {
// 						fmt.Println("error on FetchFilesAndZip:", err)
// 						RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 							EntityNumber:      entity.EntityNumber,
// 							EntityId:          order.EntityId,
// 							Indicator:         "-23",
// 							Error:             err.Error(),
// 							IsFinished:        false,
// 							Success:           "NO",
// 							AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 							AuctionStatusName: response.Orders[0].OrderStatus,
// 						})
// 						return &empty.Empty{}, err
// 					}

// 					fmt.Println("  ===> zipFileLocation :", zipFileLocation)

// 					minioLocation, err := helper.UploadToMinio(grpc.cfg, zipFileLocation, fmt.Sprintf("winner_docs_%d", req.OrderId))
// 					if err != nil {
// 						RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 							EntityNumber:      entity.EntityNumber,
// 							EntityId:          order.EntityId,
// 							Indicator:         "-24",
// 							Error:             err.Error(),
// 							IsFinished:        false,
// 							Success:           "NO",
// 							AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 							AuctionStatusName: response.Orders[0].OrderStatus,
// 						})
// 						return &empty.Empty{}, err
// 					}

// 					if minioLocation != "" {
// 						_, err = grpc.storage.FileTrack().Create(&is.CreateFileTrack{
// 							EntityId:   order.EntityId,
// 							PropertyId: ek_variables.AuctionContractAndprotocolZip,
// 							FileName:   minioLocation,
// 							BucketName: "files",
// 							ToDelete:   false,
// 							FileNameId: uuid.Nil.String(),
// 						})
// 						if err != nil {
// 							RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 								EntityNumber:      entity.EntityNumber,
// 								EntityId:          order.EntityId,
// 								Indicator:         "-32",
// 								Error:             err.Error(),
// 								IsFinished:        false,
// 								Success:           "NO",
// 								AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 								AuctionStatusName: response.Orders[0].OrderStatus,
// 							})
// 							return &empty.Empty{}, err
// 						}

// 					}

// 					err = os.RemoveAll(fmt.Sprintf("protocol_bayonnoma/%d", req.OrderId))
// 					if err != nil {
// 						RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 							EntityNumber:      entity.EntityNumber,
// 							EntityId:          order.EntityId,
// 							Indicator:         "-25",
// 							Error:             err.Error(),
// 							IsFinished:        false,
// 							Success:           "NO",
// 							AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 							AuctionStatusName: response.Orders[0].OrderStatus,
// 						})
// 						return &empty.Empty{}, err
// 					}

// 					entity, err := grpc.clients.EntityService().Get(context.Background(), &entity_service.ASGetRequest{
// 						Id: order.EntityId,
// 					})

// 					if err != nil {
// 						RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 							EntityNumber:      entity.EntityNumber,
// 							EntityId:          order.EntityId,
// 							Indicator:         "-21",
// 							Error:             err.Error(),
// 							IsFinished:        false,
// 							Success:           "NO",
// 							AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 							AuctionStatusName: response.Orders[0].OrderStatus,
// 						})
// 						return &empty.Empty{}, err
// 					}

// 					fmt.Println(" ... making davreestr request ...")
// 					// making DavAktiv request
// 					lotNumber := ""
// 					props := []*is.EntityProp{}
// 					for _, prop := range entity.EntityProperties {
// 						props = append(props, &is.EntityProp{
// 							PropertyId: prop.Property.Id,
// 							Value:      prop.Value,
// 							ValueLabel: prop.ValueLabel,
// 						})
// 						// Lot raqami
// 						if prop.Property.Id == "61a77a0676a18e5480cc5a4a" {
// 							lotNumber = prop.Value
// 						}
// 					}

// 					if lotNumber == "" {
// 						err := errors.New("lot raqami topilmadi. Yer raqami: " + entity.EntityNumber)
// 						RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 							EntityNumber:      entity.EntityNumber,
// 							EntityId:          order.EntityId,
// 							Indicator:         "-26",
// 							Error:             err.Error(),
// 							IsFinished:        false,
// 							Success:           "NO",
// 							AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 							AuctionStatusName: response.Orders[0].OrderStatus,
// 						})
// 						return &empty.Empty{}, err
// 					}

// 					if minioLocation != "" {
// 						_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(),
// 							&entity_service.AddPropertiesReq{
// 								Id: order.EntityId,
// 								Properties: []*entity_service.EntityProperty{
// 									{
// 										PropertyId: ek_variables.AuctionContractAndprotocolZip,
// 										Value:      minioLocation,
// 										ValueLabel: "Auksion shartnoma va bayonnamasi zip fayli",
// 									},
// 								},
// 							},
// 						)
// 						if err != nil {
// 							RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 								EntityNumber:      entity.EntityNumber,
// 								EntityId:          order.EntityId,
// 								Indicator:         "-33",
// 								Error:             err.Error(),
// 								IsFinished:        false,
// 								Success:           "NO",
// 								AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 								AuctionStatusName: response.Orders[0].OrderStatus,
// 							})
// 							return &empty.Empty{}, err
// 						}
// 					}

// 					minioFileUrl := grpc.cfg.MinioDomain + "/" + grpc.cfg.BucketName + "/" + minioLocation
// 					if !strings.Contains(grpc.cfg.MinioDomain, "http") {
// 						minioFileUrl = "https://" + grpc.cfg.MinioDomain + "/" + grpc.cfg.BucketName + "/" + minioLocation
// 					}

// 					result := grpc.storage.ThirdParty().MakeDavreestRequest(&is.DavreestrRequest{
// 						RequestTo:          "IJARA_AUKSION",
// 						Properties:         props,
// 						CitySoato:          int64(entity.City.Soato),
// 						RegionSoato:        int64(entity.Region.Soato),
// 						Address:            entity.Address,
// 						EntityId:           entity.Id,
// 						EntityNumber:       entity.EntityNumber,
// 						DavaktivFileLink:   minioFileUrl,
// 						DavaktivFileDate:   result.AuctionDate,
// 						DavaktivFileNumber: lotNumber,
// 					})

// 					if !result["success"].(bool) {
// 						InformWithLog(grpc.storage, "E", result["point"].(string), result["rpc-request"].(*is.DavreestrRequest), result["request"], result["response"], result["error"].(error))
// 						writeActionHistory = false
// 					} else {
// 						InformWithLog(grpc.storage, "R", result["point"].(string), result["rpc-request"].(*is.DavreestrRequest), result["request"], result["response"], nil)
// 						_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 							Id: order.EntityId,
// 							Properties: []*entity_service.EntityProperty{
// 								{
// 									PropertyId: "62df8d3729523941f03bd4be",
// 									Value:      "1",
// 								},
// 							},
// 						})

// 						if err != nil {
// 							RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 								EntityNumber:      entity.EntityNumber,
// 								EntityId:          order.EntityId,
// 								Indicator:         "-28",
// 								Error:             err.Error(),
// 								IsFinished:        false,
// 								Success:           "HALF",
// 								AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 								AuctionStatusName: response.Orders[0].OrderStatus,
// 								OldStatus:         entity.Status.Name,
// 								NewStatus:         statuses[statusSoldInAuction],
// 							})
// 							return &empty.Empty{}, err
// 						}
// 					}
// 					fmt.Println("...  successfully made davreestr request ...")
// 				}

// 				if writeActionHistory {
// 					_, err = grpc.clients.ActionHistoryService().Create(context.Background(), &user_service.ActionHistory{
// 						Id:             primitive.NewObjectID().Hex(),
// 						EntityId:       order.EntityId,
// 						Action:         recordMessage,
// 						EntityName:     entity.EntityNumber,
// 						UserId:         "000000000000000000000000",
// 						OrganizationId: "000000000000000000000000",
// 						StatusId:       statusSoldInAuction,
// 					})

// 					if err != nil {
// 						RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 							EntityNumber:      entity.EntityNumber,
// 							EntityId:          order.EntityId,
// 							Indicator:         "-20",
// 							Error:             err.Error(),
// 							IsFinished:        false,
// 							Success:           "HALF",
// 							AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 							AuctionStatusName: response.Orders[0].OrderStatus,
// 							OldStatus:         entity.Status.Name,
// 							NewStatus:         statuses[statusSoldInAuction],
// 						})
// 						return &empty.Empty{}, err
// 					}
// 				}
// 			}

// 			return &empty.Empty{}, nil

// 		// order status id = 401 or 501. Order is retracted by ekadastr staff based on some document.
// 		// Change the entity's status to Rejected with letter
// 		// order status_id=7. Order is canceled from auction
// 		case 7:
// 			new_status = statusCanceledAuction
// 			if entity.Status.Id != statusCanceledAuction {
// 				_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 					Id: order.EntityId,
// 					Properties: []*entity_service.EntityProperty{
// 						// Auksiondan bekor qilish sababi
// 						{
// 							PropertyId: "627b276f6a9eaec8edc9cf58",
// 							Value:      fmt.Sprint(response.Orders[0].LastDescription),
// 						},
// 					},
// 					StatusId: statusCanceledAuction, //statusCanceled
// 				})
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-14",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "NO",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 					})
// 					return &empty.Empty{}, err
// 				}
// 				_, err = grpc.clients.ActionHistoryService().Create(context.Background(), &user_service.ActionHistory{
// 					Id:             primitive.NewObjectID().Hex(),
// 					EntityId:       order.EntityId,
// 					Action:         "Auksiondagi yer uchastkasi savdosi bekor qilindi",
// 					EntityName:     "entity",
// 					StatusId:       statusCanceledAuction,
// 					UserId:         "000000000000000000000000",
// 					OrganizationId: "000000000000000000000000",
// 				})
// 				if err != nil {
// 					RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 						EntityNumber:      entity.EntityNumber,
// 						EntityId:          order.EntityId,
// 						Indicator:         "-15",
// 						Error:             err.Error(),
// 						IsFinished:        false,
// 						Success:           "HALF",
// 						AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 						AuctionStatusName: response.Orders[0].OrderStatus,
// 						OldStatus:         entity.Status.Name,
// 						NewStatus:         statuses[statusCanceledAuction],
// 					})
// 					return &empty.Empty{}, err
// 				}
// 			}
// 		case 401, 501:
// 			new_status = statusRevertedAuction
// 			_, err := grpc.clients.EntityService().AddPropertiesToEntity(context.Background(), &entity_service.AddPropertiesReq{
// 				Id:       order.EntityId,
// 				StatusId: statusRevertedAuction, // Auksiondan qaytarilgan
// 			})

// 			if err != nil {
// 				RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 					EntityNumber:      entity.EntityNumber,
// 					EntityId:          order.EntityId,
// 					Indicator:         "-12",
// 					Error:             err.Error(),
// 					IsFinished:        false,
// 					Success:           "NO",
// 					AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 					AuctionStatusName: response.Orders[0].OrderStatus,
// 				})
// 				return &empty.Empty{}, err
// 			}
// 			_, err = grpc.clients.ActionHistoryService().Create(context.Background(), &user_service.ActionHistory{
// 				Id:             primitive.NewObjectID().Hex(),
// 				EntityId:       order.EntityId,
// 				Action:         "Auksiondan qaytarildi",
// 				EntityName:     "entity",
// 				StatusId:       statusRevertedAuction,
// 				UserId:         "000000000000000000000000",
// 				OrganizationId: "000000000000000000000000",
// 			})
// 			if err != nil {
// 				RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 					EntityNumber:      entity.EntityNumber,
// 					EntityId:          order.EntityId,
// 					Indicator:         "-13",
// 					Error:             err.Error(),
// 					IsFinished:        false,
// 					Success:           "HALF",
// 					AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 					AuctionStatusName: response.Orders[0].OrderStatus,
// 					OldStatus:         entity.Status.Name,
// 					NewStatus:         statuses[statusRevertedAuction],
// 				})
// 				return &empty.Empty{}, err
// 			}
// 		case 301:
// 			_, err = grpc.clients.ActionHistoryService().Create(context.Background(), &user_service.ActionHistory{
// 				Id:             primitive.NewObjectID().Hex(),
// 				EntityId:       order.EntityId,
// 				Action:         "Auksionda shartnoma imzolanmadi",
// 				EntityName:     entity.EntityNumber,
// 				StatusId:       entity.Status.Id,
// 				UserId:         "000000000000000000000000",
// 				OrganizationId: "000000000000000000000000",
// 			})
// 			if err != nil {
// 				RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 					EntityNumber:      entity.EntityNumber,
// 					EntityId:          order.EntityId,
// 					Indicator:         "-31",
// 					Error:             err.Error(),
// 					IsFinished:        false,
// 					Success:           "HALF",
// 					AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 					AuctionStatusName: response.Orders[0].OrderStatus,
// 					OldStatus:         entity.Status.Name,
// 					NewStatus:         statuses[statusRevertedAuction],
// 				})
// 				return &empty.Empty{}, err
// 			}
// 		}
// 	case 50:
// 	case 37:
// 	}

// 	newStatus := "----"
// 	if _, ok := statuses[new_status]; ok {
// 		newStatus = statuses[new_status]
// 	}
// 	RecordPushFromAuction(grpc.storage, req, v_is.AuctionPushLog{
// 		EntityNumber:      entity.EntityNumber,
// 		EntityId:          order.EntityId,
// 		Indicator:         "-S1",
// 		IsFinished:        true,
// 		Success:           "YES",
// 		AuctionStatusCode: response.Orders[0].OrderStatusesID,
// 		AuctionStatusName: response.Orders[0].OrderStatus,
// 		OldStatus:         entity.Status.Name,
// 		NewStatus:         newStatus,
// 	})

// 	return &empty.Empty{}, nil
// }

// func RecordPushFromAuction(store storage.StorageI, aucReq *is.PushFromAuctionRequest, req v_is.AuctionPushLog) error {
// 	beautifulJsonByte, err := json.MarshalIndent(aucReq, "", "  ")
// 	body := ""
// 	if err != nil {
// 		body = fmt.Sprintf("%v", aucReq)
// 	} else {
// 		body = string(beautifulJsonByte)
// 	}

// 	err = store.ActionLog().RecordPushFromAuction(v_is.AuctionPushLog{
// 		OrderId:    aucReq.OrderId,
// 		CategoryId: aucReq.CategoryId,
// 		Protocol:   aucReq.Protocol,
// 		Url:        aucReq.Url,

// 		EntityId:          req.EntityId,
// 		EntityNumber:      req.EntityNumber,
// 		Indicator:         "POINT" + req.Indicator,
// 		OldStatus:         req.OldStatus,
// 		NewStatus:         req.NewStatus,
// 		AuctionStatusName: req.AuctionStatusName,
// 		AuctionStatusCode: req.AuctionStatusCode,
// 		IsFinished:        req.IsFinished,
// 		Success:           req.Success,
// 		PushReqRody:       body,
// 		Error:             req.Error,
// 	})
// 	if req.Error != "" {
// 		fmt.Println("ERROR IN PUSH: ", req.Error, "point: ", req.Indicator)
// 	}
// 	return err
// }

// func InformWithLog(strg storage.StorageI, type_ string, indicator string, req *is.DavreestrRequest, reqBody interface{}, respBody interface{}, err error) {
// 	if err != nil {
// 		fmt.Println("InformWithLog >>> ", err)
// 	}

// 	var errOnRecord error
// 	if type_ == "E" {
// 		errOnRecord = strg.FunctionLogs().CreateIntegrationErrorLog(&is.CreateErrorLogRequest{
// 			Indicator:   "DR_ERR_POINT-" + indicator,
// 			Target:      req.EntityNumber,
// 			OtherTarget: fmt.Sprintf("%d", req.OrderId),
// 			TargetId:    req.EntityId,
// 			EndpointUrl: ek_variables.DavReestrUrl,
// 			ResBody:     util.JSONStringify(err),
// 			ReqBody:     util.JSONStringify(reqBody),
// 		})
// 	} else if type_ == "R" {
// 		errOnRecord = strg.FunctionLogs().CreateIntegrationInfoLog(&is.CreateInfoLogRequest{
// 			Indicator:   "DR_REC_POINT-" + indicator,
// 			Target:      req.EntityNumber,
// 			OtherTarget: fmt.Sprintf("%d", req.OrderId),
// 			TargetId:    req.EntityId,
// 			EndpointUrl: ek_variables.DavReestrUrl,
// 			ResBody:     util.JSONStringify(respBody),
// 			ReqBody:     util.JSONStringify(reqBody),
// 		})
// 	}

// 	if errOnRecord != nil {
// 		fmt.Println("-------------------------------- > could not record (third party)1 <-------------------------------", errOnRecord)
// 	}
// }
