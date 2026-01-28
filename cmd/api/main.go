package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/middleware"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/routes"
)

func main() {
	// Entry point for the API server

	// Load environment variables
	env := config.LoadEnv()

	// Load DB URL from environment
	dsn := env.DatabaseUrl
	if dsn == "" {
		log.Fatal("DB_URL is not set")
	}

	// Connect to the database
	db, err := config.Connect(dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	_ = db // just to show it’s connected

	// Connect to Redis
	config.ConnectRedis()

	var router *gin.Engine = gin.Default()
	//router := gin.Default()

	router.SetTrustedProxies(nil)
	router.Use(middleware.CORSMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		fmt.Println("go working")
		ctx.JSON(200, gin.H{
			"message": "go working",
			"status":  "success",
		})
	})

	// Call SetRoutes to register all API routes
	routes.SetRoutes(router)

	// Start server
	port := env.Port
	if port == "" {
		port = "8080" // default
	}
	router.Run(":" + port)

}
