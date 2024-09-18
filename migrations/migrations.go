package migrations

import (
	"ecommerce-api/config"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/repo/driver"
	"fmt"

	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func Migration(operation string) {
	// init the env config
	cfg, err := config.LoadConfig(consts.AppName)
	if err != nil {
		panic(err)
	}
	// logrus init
	log := logrus.New()

	// database connection
	pgsqlDB, err := driver.ConnectDB(cfg.Db)
	if err != nil {
		log.Fatalf("unable to connect to the database")
		return
	}
	driverDB, err := postgres.WithInstance(pgsqlDB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	// migration instance creation
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations/scripts",
		consts.DatabaseType, driverDB)
	if err != nil {
		panic(err)
	}

	switch operation {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			panic(err)
		} else if err == migrate.ErrNoChange {
			fmt.Println("No migration to apply")
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			panic(err)
		} else if err == migrate.ErrNoChange {
			fmt.Println("No migration to rollback")
		}
	default:
		fmt.Printf("Unsupported migration operation: %s\n", operation)
		return
	}

	log.Printf("Migration %v completed...", operation)
}
