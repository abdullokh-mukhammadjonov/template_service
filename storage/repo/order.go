package repo

import (
	"gitlab.udevs.io/ekadastr/ek_integration_service/genproto/content_service"
)

type OrderI interface {
	Create(req map[string]interface{}, lawType string) error
	Get(*content_service.GetHandbooksRequest) (*content_service.GetHandbooksResponse, error)
	GetOne(*content_service.GetOneRequest) (*content_service.GetOneHandbookResponse, error)
}
