package storage

import (
	"github.com/abdullokh-mukhammadjonov/template_service/storage/postgres"
	"github.com/abdullokh-mukhammadjonov/template_service/storage/repo"
	"github.com/jmoiron/sqlx"
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
