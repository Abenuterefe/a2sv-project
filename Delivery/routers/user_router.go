package routers

import (
	"github.com/Abenuterefe/a2sv-project/delivery/controllers"
	"github.com/Abenuterefe/a2sv-project/infrastructure/auth"
	"github.com/Abenuterefe/a2sv-project/infrastructure/middlewares"
	"github.com/Abenuterefe/a2sv-project/repository"
	"github.com/Abenuterefe/a2sv-project/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(r *gin.Engine, mongoClient *mongo.Client) {
	// Initialzie MongoDB database
	db := mongoClient.Database("g6_starter_projectDb")

	// Set up repository, usecase, password service, controller and other services
	PasswordService := auth.NewBcryptPasswordService()
	jwtService := auth.NewJWTService()

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUsecase(userRepo, PasswordService, jwtService)
	authCtrl := controllers.NewAuthController(userUseCase)

	// Setup /auth group routes
	authGroup := r.Group("/auth")
	authGroup.POST("/register", authCtrl.Regiser)
	authGroup.POST("/login", authCtrl.Login)
	authGroup.POST("/refresh", authCtrl.Refresh)

	// Setup protected routes
	protected := r.Group("/user")
	protected.Use(middlewares.AuthMiddleware(jwtService))
	//-----Regular user only------- //
	protected.GET("/profile", middlewares.UserOnlyMiddleware(), authCtrl.Profile)
	// ---- Admin dashboard (admin only)-----  //
	protected.GET("/admin", middlewares.AdminOnlyMiddleware(), authCtrl.AdminDashboard)
}
