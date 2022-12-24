package repo

import (
	"gitlab.udevs.io/ekadastr/ek_integration_service/genproto/content_service"
)

type OrderI interface {
	Create(req map[string]interface{}, lawType string) error
	Get(content_service.) (content_service.ExcelReportResponse, error)
	GetOne(content_service.ExcelReportRequest) (content_service.ExcelReportResponse, error)
}
