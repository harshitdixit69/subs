package main

import (
	"log"
	"os"

	"github.com/fliptable-io/subscription-service/src/data"
	"github.com/fliptable-io/subscription-service/src/repos"
	"github.com/fliptable-io/subscription-service/src/routers"
	"github.com/fliptable-io/subscription-service/src/services"
	loggy "github.com/fliptable-io/subscription-service/src/utils/logging"
	"github.com/joho/godotenv"
)

func main() {
	//HttpRepo := &repos.HttpRepo{}

	//loads values from .env into the system
	if err := godotenv.Load(".main.env"); err != nil {
		log.Print("No .env file found")
	}

	pg := repos.NewPostgresRepo()

	// services
	subService := services.SubscriptionService{Pg: pg}

	// routers
	subRouter := routers.SubscriptionRouter{SubscriptionService: subService}
	baseRouter := routers.NewBaseRouter(subRouter)

	port := os.Getenv("SERVER_PORT")
	loggy.Info(data.Saitama(port))
	err := baseRouter.Run(":" + port)
	if err != nil {	
		log.Fatal(err)
	}
}
