package routers

import (
	"github.com/Abenuterefe/a2sv-project/delivery/controllers"
	"github.com/Abenuterefe/a2sv-project/infrastructure/auth"
	"github.com/Abenuterefe/a2sv-project/infrastructure/middlewares"
	"github.com/Abenuterefe/a2sv-project/infrastructure/storage"
	"github.com/Abenuterefe/a2sv-project/repository"
	"github.com/Abenuterefe/a2sv-project/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProfileRoutes(r *gin.Engine, mongoClient *mongo.Client) {
	// Initialize dependencies
	db := mongoClient.Database("g6_starter_projectDb")
	jwtService := auth.NewJWTService()
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)

	// Create file storage service (local storage)
	fileStorage := storage.NewLocalFileStorage("uploads/profile_pictures")

	// Pass all required dependencies to usecase
	profileUsecase := usecase.NewProfileUsecase(userRepo, profileRepo, fileStorage)
	profileController := controllers.NewProfileController(profileUsecase)

	// Protected routes for /user/profile
	profileGroup := r.Group("/user")
	profileGroup.Use(middlewares.AuthMiddleware(jwtService))
	{
		// Update text fields (username, bio, etc.)
		profileGroup.PUT("/profile", profileController.UpdateProfile)

		// Get full profile (email, username, bio, picture)
		profileGroup.GET("/profile/me", profileController.GetProfile)

		// Upload profile picture separately
		profileGroup.POST("/profile/picture", profileController.UploadProfilePicture)
	}
}
