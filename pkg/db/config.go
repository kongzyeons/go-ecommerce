package db

import (
	"app-ecommerce/config"
	"database/sql"
	"errors"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBX struct {
	sqlx.DB
}

var dbx *DBX
var dbxOnce sync.Once

// InitDatabase Function for init database
func InitDatabase(connectionString string) {
	dbxOnce.Do(func() {
		dbs, err := sqlx.Connect("postgres", connectionString)
		if err != nil {
			log.Panic(err)
		}

		dbx = &DBX{
			DB: *dbs,
		}

		config := config.GetConfig()
		if config.PostgresDB.MaxOpenConn > 0 {
			dbx.SetMaxOpenConns(config.PostgresDB.MaxOpenConn) // The default is 0 (unlimited)
		}

		if config.PostgresDB.MaxIdleConn > 0 {
			dbx.SetMaxIdleConns(config.PostgresDB.MaxIdleConn) // defaultMaxIdleConns = 2
		}

		if config.PostgresDB.ConnMaxLifeTimeTTL != nil {
			dbx.SetConnMaxLifetime(*config.PostgresDB.ConnMaxLifeTimeTTL) // 0, connections are reused forever.
		}
	})
}

// UnInitDatabase cleanup database
func UnInitDatabase() {
	if dbx != nil {
		dbx.DB.Close()
	}
}

// IsSQLReallyError check if really error
func IsSQLReallyError(err error) bool {
	return err != nil && !errors.Is(err, sql.ErrNoRows)
}
