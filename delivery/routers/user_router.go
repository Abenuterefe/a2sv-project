package routers

import (
	"github.com/Abenuterefe/a2sv-project/delivery/controllers"
	"github.com/Abenuterefe/a2sv-project/infrastructure/auth"
	"github.com/Abenuterefe/a2sv-project/infrastructure/mail"

	"github.com/Abenuterefe/a2sv-project/infrastructure/middlewares"
	"github.com/Abenuterefe/a2sv-project/repository"
	"github.com/Abenuterefe/a2sv-project/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(r *gin.Engine, mongoClient *mongo.Client) {
	// Initialzie MongoDB database
	db := mongoClient.Database("g6_starter_projectDb")

	// Initialize services
	PasswordService := auth.NewBcryptPasswordService()
	jwtService := auth.NewJWTService()
	mailService := mail.NewMailService()
	oauthService := auth.NewGoogleOAuthService()
	secureToken := auth.NewSecureTokenGenerator()

	// Initialize repository, usecase, controller
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUsecase(userRepo, PasswordService, jwtService, mailService, oauthService, secureToken)
	authCtrl := controllers.NewAuthController(userUseCase, oauthService)

	// ----------------------
	// Auth Routes (/auth)
	// ----------------------
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authCtrl.Regiser)
		authGroup.POST("/login", authCtrl.Login)
		authGroup.POST("/refresh", authCtrl.Refresh)
		authGroup.GET("/verify", authCtrl.VerifyEmail)
		authGroup.POST("/resend-verification", authCtrl.ResendVerification)
		//login with google
		authGroup.GET("/google/login", authCtrl.GoogleLogin)
		authGroup.GET("/google/callback", authCtrl.GoogleCallback)
		//forget and reset pwd
		authGroup.POST("/forgot-password", authCtrl.ForgotPassword)
		authGroup.POST("/reset-password", authCtrl.ResetPassword)
	}

	// ----------------------
	// Protected Routes (/user)
	// ----------------------
	protected := r.Group("/user")
	protected.Use(middlewares.AuthMiddleware(jwtService))

	//-----Regular user only------- //
	{
		protected.GET("/profile", middlewares.UserOnlyMiddleware(), authCtrl.Profile)
	}

	// ------Logout------//
	protected.POST("/logout", authCtrl.Logout)

	// ----------------------
	// Admin-only Routes (/user/admin)
	// ----------------------
	adminGroup := protected.Group("/admin")
	adminGroup.Use(middlewares.AdminOnlyMiddleware())
	{
		adminGroup.GET("", authCtrl.AdminDashboard)
		adminGroup.PUT("/promote/:id", authCtrl.PromoteUser)
		adminGroup.PUT("/demote/:id", authCtrl.DemoteUser)
	}
}
