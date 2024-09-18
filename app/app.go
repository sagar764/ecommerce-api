package app

import (
	"context"
	"ecommerce-api/config"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/controllers"
	"ecommerce-api/internal/entities"
	"ecommerce-api/internal/middlewares"
	"ecommerce-api/internal/repo"
	"ecommerce-api/internal/repo/driver"
	"ecommerce-api/internal/usecases"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "ecommerce-api/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// method run
// env configuration
// logrus, zap
// use case intia
// repo initalization
// controller init

func Run() {
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
		log.Fatalf("unable to connect the database : %v", err)
		return
	}

	// here initalizing the router
	router := initRouter()
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	api := router.Group("/api")
	// middleware initialization
	m := middlewares.NewMiddlewares(cfg)
	api.Use(m.ApiVersioning())

	// complete ecommerce related initialization
	{

		// repo initialization
		ecommerceRepo := repo.NewEcommerceRepo(pgsqlDB)

		// initilizing usecases
		ecommerceUseCases := usecases.NewEcommerceUseCases(ecommerceRepo)

		// initalizin controllers
		ecommerceControllers := controllers.NewEcommerceController(api, m, ecommerceUseCases)

		// init the routes
		ecommerceControllers.InitRoutes()
	}

	// runn the app
	launch(cfg, router)
}

func initRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// CORS
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// },
		MaxAge: 12 * time.Hour,
	}))

	// common middlewares should be added here

	return router
}

// launch
func launch(cfg *entities.EnvConfig, router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	fmt.Println("Server listening in...", cfg.Port)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
