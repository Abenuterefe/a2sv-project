package main

import (
	"log"
	"os"

	"github.com/Abenuterefe/a2sv-project/delivery/routers"
	"github.com/Abenuterefe/a2sv-project/infrastructure/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file is not find")
	}

	// =========================================  //
	// CONNECT TO MONGODB DATABASE
	mongoClient, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	// CALL AND CREATE ROUTER
	r := gin.Default()
	routers.UserRoutes(r, mongoClient)
	routers.AiRoutes(r)

	// RUN SERVER
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 👇starts the actual server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
