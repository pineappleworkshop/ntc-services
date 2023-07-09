package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"ntc-services/config"
	"ntc-services/handlers"
	"ntc-services/services"
	"ntc-services/stores"
	"os"
	"strconv"
)

func init() {
	config.InitConf()
	stores.InitDbs()
	confLogLevel, err := config.GetLogLevel()
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	logLevel, err := log.ParseLevel(*confLogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if err := services.StartServices(); err != nil {
		panic(err)
	}

	initPublicRoutes(e)
	go log.Fatal(e.Start(":" + strconv.Itoa(config.PORT)))
}

func initPublicRoutes(e *echo.Echo) {
	// health
	e.GET("/health", handlers.HealthHandler)

}
