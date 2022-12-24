package storage

import (
	"github.com/jmoiron/sqlx"
	"gitlab.udevs.io/ekadastr/ek_integration_service/storage/postgres"
	"gitlab.udevs.io/ekadastr/ek_integration_service/storage/repo"
)

type StorageI interface {
	Order() repo.OrderI
}

type storagePostgres struct {
	db        *sqlx.DB
	orderRepo repo.OrderI
}

func NewStoragePostgres(db *sqlx.DB) StorageI {
	return &storagePostgres{
		orderRepo: postgres.NewOrderRepo(db),
	}
}

func (s *storagePostgres) Order() repo.OrderI {
	return s.orderRepo
}
