package main

import (
	"log"
	"os"

	"github.com/Abenuterefe/a2sv-project/delivery/routers"
	"github.com/Abenuterefe/a2sv-project/infrastructure/database"
	"github.com/Abenuterefe/a2sv-project/infrastructure/ai"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file is not found")
	}
	ai.Setup()

	// CONNECT TO MONGODB DATABASE
	mongoClient, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	// CREATE ROUTER
	r := gin.Default()

	// ✅ Serve static files (images) from /uploads
	r.Static("/uploads", "./uploads")

	// ROUTES
	routers.BlogRoutes(r, mongoClient)
	routers.UserRoutes(r, mongoClient)
	routers.ProfileRoutes(r, mongoClient)
	routers.AiRoutes(r)
	routers.CommentRoutes(r, mongoClient)
	routers.BlogInteractionRoutes(r, mongoClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// START SERVER
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
