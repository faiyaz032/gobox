package repository

import (
	db "github.com/faiyaz032/gobox/internal/infra/db/sqlc"
)

type BoxRepository struct {
	queries *db.Queries
}

func NewBoxRepository(queries *db.Queries) *BoxRepository {
	return &BoxRepository{
		queries: queries,
	}
}
