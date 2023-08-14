package main

import (
	"ntc-services/config"
	"ntc-services/handlers"
	"ntc-services/services"
	"ntc-services/stores"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
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

	// statspool
	e.GET("/statspool", handlers.GetStatsPool)

	// inscriptions
	e.GET("/addresses/:addr/inscriptions", handlers.GetInscriptions)
	e.GET("/addresses/:addr/brc20s", handlers.GetBRC20s)

	// wallets
	e.POST("/wallets", handlers.PostWallets)
	//e.GET("/wallets/:id/inscriptions", handlers.GetInscriptions) // TODO: implement
	//e.GET("/wallets/:id/brc20s", handlers.GetBRC20s)			   // TODO: implement

	// trades
	e.POST("/trades", handlers.PostTrades)
	e.POST("/trades/:id/maker", handlers.PostMakerByTradeID)
	e.GET("/trades", handlers.GetTrades)
	e.POST("/trades/:id/offers", handlers.PostOfferByTradeID)
	e.GET("/trades/:id/offers", handlers.GetOffersByTradeID)
	e.POST("/trades/:id/orders/accept", handlers.PostAcceptOfferByTradeID)
	e.POST("/trades/:id/submit", handlers.PostSubmitTradeByID)

	// ordex inscription testing api
	e.GET("/ordex/inscription/:id", handlers.OrdexHandler)

	// experiments
	e.GET("/experiments/from-unsigned-tx", handlers.PSBTFromUnsignedTx)
	e.GET("/experiments/psbt", handlers.GeneratePSBT)
	e.GET("/experiments/utxos", handlers.UTXOs)
}
