package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/anssuy/code-colosseum/backend/internal/auth"
	"github.com/anssuy/code-colosseum/backend/internal/db"
	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	auth.Init()

	pool, err := db.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := dbgen.New(pool)

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	tokenManager := auth.NewTokenManager(jwtSecret)
	authHandler := auth.NewHandler(queries, tokenManager)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Code Colosseum API"})
	})

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.POST("/logout", authHandler.Logout)

		authRoutes.GET("/me", auth.Middleware(tokenManager), authHandler.Me)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
