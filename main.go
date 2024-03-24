package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/api/routes"
	"github.com/hodukihugi/winglets-api/configs"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	router := SetupRouter()
	router.Run(os.Getenv("PORT"))
}

func SetupRouter() *gin.Engine {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := configs.Connection()
	router := gin.Default()

	/**
	@description Init All Route
	*/

	routes.InitUserRoutes(db, router)

	return router
}
