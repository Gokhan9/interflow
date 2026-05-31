package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(connString string) error {
	var err error
	DB, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		return err
	}

	// bağlantı testi için ping atma işlemi..
	return DB.Ping(context.Background())
}
