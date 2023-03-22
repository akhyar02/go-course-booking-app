package dbrepo

import (
	"database/sql"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/repository"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		App: a,
		DB:  conn,
	}
}
