package main

import (
	"ecommerce-api/app"
	"ecommerce-api/migrations"

	"flag"

	"github.com/sirupsen/logrus"
)

var (
	runserver = flag.Bool("runserver", false, "This is a string argument for running server")
	migration = flag.Bool("migration", false, "This is a string argument for running migration")
	up        = flag.Bool("up", false, "This is a string argument for running migration up")
	down      = flag.Bool("down", false, "This is a string argument for running migration down")
)

// @title E-commerce API
// @version 1.0
// @description This is an e-commerce API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// logrus init
	log := logrus.New()

	flag.Parse()

	if !*runserver && !*migration {
		log.Fatalf("Please specify the file you want to execute")
	}
	if *runserver {
		app.Run()
	}
	if *migration {
		if !*down && !*up {
			log.Fatal("Please specify the migration type")
		}
		if *up {
			migrations.Migration("up")
		}
		if *down {
			migrations.Migration("down")
		}
	}

}
