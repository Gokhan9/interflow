package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool" //PostgreSQL bağlantı havuzu yönetimini sağlayan kütüphane.
)

/*
DB-postgresql connection pool tutar
* *pgxpool.pool connection pool. Globaldir, her yerden erişim.
connString string == "postgres://user:password@localhost:5432/mydb" gibi bir bağlantı adresini tutar.
pgxpool.New, yeni connection pool oluşturur, gerçek db bağlantıları gerekince açılır. Oluşturulan "pool" global DB değişkeninde tutulur.
*/
var DB *pgxpool.Pool

func InitDB(connString string) error {

	var err error //hata bilgisi

	DB, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		return err
	}

	return DB.Ping(context.Background())
}
