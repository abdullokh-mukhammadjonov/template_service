package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	cs "gitlab.udevs.io/ekadastr/ek_integration_service/genproto/content_service"
	"gitlab.udevs.io/ekadastr/ek_integration_service/storage/repo"
)

type orderRepo struct {
	db *sqlx.DB
}

// NewOrderRepo ...
func NewOrderRepo(db *sqlx.DB) repo.OrderI {
	return &orderRepo{db: db}
}

func (or *orderRepo) Create(order map[string]interface{}, lawType string) error {
	if lawType == "1" {
		order["law_type"] = "mulk huquqi"
	} else if lawType == "0" {
		order["law_type"] = "ijara huquqi"
	}
	logCreateQuery := `
		INSERT INTO
			orders(
				entity_id,
				order_id,
				law_type,
				language,
				soato,
				name,
				price,
				address,
				category,
				closed,
				prop_set,
				additional_info,
				lat,
				lng,
				created_at
			)
		VALUES (
			:entity_id,
			:order_id,
			:law_type,
			:language,
			:soato,
			:name,
			:price,
			:address,
			:category,
			:closed,
			:prop_set,
			:additional_info,
			:lat,
			:lng,
			now()
		)
	`

	_, err := or.db.NamedExec(logCreateQuery, order)
	if err != nil {
		return err
	}

	return nil
}

func (or *orderRepo) Update(order map[string]interface{}, lawType string) error {
	if lawType == "1" {
		order["law_type"] = "mulk huquqi"
	} else if lawType == "0" {
		order["law_type"] = "ijara huquqi"
	}
	updateQuery := `
		UPDATE	
			orders
		SET
			order_id = :order_id,
			language = :language,
			soato = :soato,
			name = :name,
			price = :price,
			address = :address,
			category = :category,
			closed = :closed,
			prop_set = :prop_set,
			additional_info = :additional_info,
			lat = :lat,
			lng = :lng,
			law_type = :law_type
		WHERE entity_id = :entity_id
	`

	_, err := or.db.NamedExec(updateQuery, order)
	if err != nil {
		return err
	}

	return nil
}

func (or *orderRepo) Get(req *cs.GetHandbooksRequest) (res *cs.GetHandbooksResponse, err error) {
	fmt.Println("storage.GetOrder.req.entityID:")
	getQuery := `
		SELECT
			file_link,
		FROM 
			orders
		WHERE
			file_link=$1
	`

	row := or.db.QueryRow(getQuery, &req)
	err = row.Scan(
		&res,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (or *orderRepo) GetOne(req *cs.GetOneRequest) (res *cs.GetOneHandbookResponse, err error) {
	fmt.Println("storage.GetOrder.req.entityID:")
	getQuery := `
		SELECT
			file_link,
		FROM 
			orders
		WHERE
			file_link=$1
	`

	row := or.db.QueryRow(getQuery, &req)
	err = row.Scan(
		&res,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}
