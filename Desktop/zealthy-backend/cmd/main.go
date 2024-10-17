package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tbui1996/zealthy-backend/internal/adapters/handlers"
	"github.com/tbui1996/zealthy-backend/internal/adapters/repositories"
	"github.com/tbui1996/zealthy-backend/internal/core/services"
	"github.com/tbui1996/zealthy-backend/pkg/database"
)

func main() {
	// Database connection
	db, err := database.NewPostgresConnection("postgres://zealthy_user:your_secure_password@localhost:5432/zealthy?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repositories and services
	userRepo := repositories.NewPostgresUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	// Custom CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true // Be cautious with this in production
		},
		MaxAge: 12 * time.Hour,
	}))

	// Add a global OPTIONS handler
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours
		c.Status(204)
	})

	// User routes
	r.POST("/user", userHandler.CreateUser)
	r.PUT("/user/:id", userHandler.UpdateUser)
	r.GET("/users/:id", userHandler.GetUser)
	r.GET("/users", userHandler.GetAllUsers)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.NoRoute(func(c *gin.Context) {
		log.Printf("No route found for %s %s", c.Request.Method, c.Request.URL)
		c.JSON(404, gin.H{"error": "Route not found"})
	})
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
