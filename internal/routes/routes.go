package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-enterprise-api/internal/config"
	"github.com/yourusername/go-enterprise-api/internal/database"
	"github.com/yourusername/go-enterprise-api/internal/handlers"
	"github.com/yourusername/go-enterprise-api/internal/middleware"
	"github.com/yourusername/go-enterprise-api/internal/repository"
	"github.com/yourusername/go-enterprise-api/internal/services"
)

// Setup configures all routes
func Setup(cfg *config.Config, db *database.Database) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS(&cfg.CORS))
	router.Use(middleware.RateLimit(cfg.RateLimit.Requests, cfg.RateLimit.Duration))

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	postRepo := repository.NewPostRepository(db.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	healthHandler := handlers.NewHealthHandler(db)

	// API version group
	api := router.Group("/api/v1")

	// Health routes (no authentication required)
	healthRoutes := api.Group("/health")
	{
		healthRoutes.GET("", healthHandler.Health)
		healthRoutes.GET("/ready", healthHandler.Ready)
		healthRoutes.GET("/live", healthHandler.Live)
	}

	// Auth routes
	authRoutes := api.Group("/auth")
	{
		// Public routes with stricter rate limiting
		publicAuth := authRoutes.Group("")
		publicAuth.Use(middleware.StrictRateLimit(10, time.Minute))
		{
			publicAuth.POST("/register", authHandler.Register)
			publicAuth.POST("/login", authHandler.Login)
			publicAuth.POST("/refresh", authHandler.RefreshTokens)
		}

		// Protected routes
		protectedAuth := authRoutes.Group("")
		protectedAuth.Use(middleware.AuthMiddleware(authService))
		{
			protectedAuth.POST("/logout", authHandler.Logout)
			protectedAuth.GET("/me", authHandler.Me)
			protectedAuth.POST("/change-password", authHandler.ChangePassword)
		}
	}

	// User routes
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(authService))
	{
		// Standard user routes
		userRoutes.GET("", userHandler.GetAll)
		userRoutes.GET("/search", userHandler.Search)
		userRoutes.GET("/:id", userHandler.GetByID)
		userRoutes.PUT("/:id", userHandler.Update)

		// Admin only routes
		adminRoutes := userRoutes.Group("")
		adminRoutes.Use(middleware.RequireAdmin())
		{
			adminRoutes.DELETE("/:id", userHandler.Delete)
			adminRoutes.PATCH("/:id/status", userHandler.UpdateStatus)
			adminRoutes.PATCH("/:id/role", userHandler.UpdateRole)
		}
	}

	// Post routes
	postRoutes := api.Group("/posts")
	{
		// Public routes (with optional auth for viewing drafts)
		postRoutes.GET("", middleware.OptionalAuthMiddleware(authService), postHandler.GetAll)
		postRoutes.GET("/search", postHandler.Search)
		postRoutes.GET("/slug/:slug", middleware.OptionalAuthMiddleware(authService), postHandler.GetBySlug)
		postRoutes.GET("/:id", middleware.OptionalAuthMiddleware(authService), postHandler.GetByID)

		// Protected routes
		protectedPosts := postRoutes.Group("")
		protectedPosts.Use(middleware.AuthMiddleware(authService))
		{
			protectedPosts.POST("", postHandler.Create)
			protectedPosts.GET("/my", postHandler.GetMyPosts)
			protectedPosts.PUT("/:id", postHandler.Update)
			protectedPosts.DELETE("/:id", postHandler.Delete)
		}
	}

	// Admin routes
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(authService))
	adminRoutes.Use(middleware.RequireAdmin())
	{
		adminRoutes.GET("/health/info", healthHandler.Info)
	}

	return router
}
