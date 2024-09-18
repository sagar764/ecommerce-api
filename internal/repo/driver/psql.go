package driver

import (
	"database/sql"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"fmt"

	// Import the "github.com/jackc/pgx/v4" package for PostgreSQL database connectivity.
	// This package provides a Go driver for PostgreSQL databases.
	_ "github.com/jackc/pgx/v4"

	// Import the "github.com/lib/pq" package for PostgreSQL database connectivity.
	// This package is another Go driver commonly used for working with PostgreSQL databases.
	_ "github.com/lib/pq"
)

// ConnectDB initializes postgres DB
func ConnectDB(cfg entities.Database) (*sql.DB, error) {
	datasource := prepareConnectionString(cfg)
	databaseType := consts.DatabaseType
	db, err := sql.Open(databaseType, datasource)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %s err: %s", datasource, err)
	}
	db.SetMaxOpenConns(cfg.MaxActive)
	db.SetMaxIdleConns(cfg.MaxIdle)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db(ping): %s err: %s", datasource, err)
	}
	return db, nil
}

func prepareConnectionString(cfg entities.Database) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=20 search_path=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DATABASE, cfg.Schema)
}
