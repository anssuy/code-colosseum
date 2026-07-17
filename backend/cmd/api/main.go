package main

import (
	"context"
	"log"
	"os"

	"github.com/anssuy/code-colosseum/backend/internal/auth"
	db "github.com/anssuy/code-colosseum/backend/internal/db/generated"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pool, err := pgxpool.New(
		context.Background(),
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	queries := db.New(pool)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	tokenManager := auth.NewTokenManager(jwtSecret)
	authHandler := auth.NewHandler(queries, tokenManager)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Code Colosseum API",
		})
	})

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)

	router.Run(":8080")
}
