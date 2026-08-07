package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"payment_verifier/config"
	"payment_verifier/controllers"
	"payment_verifier/middleware"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Key, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	config.LoadConfig()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())
	r.Use(middleware.StructuredLogger())
	r.Use(middleware.RateLimitMiddleware(20, 50)) // 20 requests/sec, burst 50

	// Health and observability routes (public)
	r.GET("/health", controllers.HealthHandler)
	r.GET("/ping", controllers.HealthHandler)
	r.GET("/metrics", controllers.MetricsHandler)

	// API routes (protected if API_KEY is set in environment)
	api := r.Group("/api/verify")
	api.Use(middleware.APIKeyAuth())
	{
		api.POST("", controllers.VerifyHandler)
		api.POST("/", controllers.VerifyHandler)
		api.POST("/batch", controllers.BatchVerifyHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Ethiopian Payment Receipt Verifier running on http://localhost:%s", port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
