package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/internal/apps/base/controllers/authentication"
	"github.com/sdq-codes/usegro-api/internal/apps/base/controllers/user"
	"github.com/sdq-codes/usegro-api/internal/apps/base/controllers/verification"
	authenticationService "github.com/sdq-codes/usegro-api/internal/apps/base/services/authentication"
	userService "github.com/sdq-codes/usegro-api/internal/apps/base/services/user"
	"github.com/sdq-codes/usegro-api/internal/router/middleware"
	"gorm.io/gorm"
)

func BaseRouter(v1 fiber.Router, db *gorm.DB, rdb *redis.Client) {
	// Authentication
	authService := authenticationService.NewAuthenticationService(db, rdb)

	authAPIGroup := v1.Group("/authentication")
	authController := authentication.NewAuthenticationController(*authService)
	authAPIGroup.Post("/email/exist", authController.Exist)
	authAPIGroup.Post("/register", authController.Register)
	authAPIGroup.Post("/login", authController.Login)
	authAPIGroup.Post("/refresh", authController.Refresh)
	authAPIGroup.Post("/logout", authController.Logout)
	authAPIGroup.Post("/forgot-password", authController.ForgotPassword)
	authAPIGroup.Post("/reset-password", authController.ResetPassword)
	authAPIGroup.Post("/email-code/request", authController.RequestEmailCode)
	authAPIGroup.Post("/email-code/verify", authController.VerifyEmailCode)

	// Google OAuth
	googleAuthAPIGroup := authAPIGroup.Group("/google")
	googleAuthController := authentication.NewGoogleAuthenticationController(authService, rdb)
	googleAuthAPIGroup.Get("/login", googleAuthController.HandleGoogleLogin)
	googleAuthAPIGroup.Get("/callback", googleAuthController.HandleGoogleCallback)

	// User
	userAPIGroup := v1.Group("/user")
	userController := user.NewUserController(*userService.NewUserService(db))
	userAPIGroup.Get("", userController.GetLoggedInUser)

	// Email verification
	verificationAPIGroup := v1.Group("/verification")
	emailVerificationAPIGroup := verificationAPIGroup.Group("/email")
	emailVerificationController := verification.NewEmailVerificationController(db, rdb)
	emailVerificationAPIGroup.Post("/", middleware.JwtVerify(), emailVerificationController.Verify)
	emailVerificationAPIGroup.Get("/resend", middleware.JwtVerify(), emailVerificationController.Resend)
}
