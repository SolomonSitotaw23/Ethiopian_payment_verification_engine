package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"payment_verifier/services"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "UP",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
	})
}

func MetricsHandler(c *gin.Context) {
	snapshot := services.Metrics.GetSnapshot()
	c.JSON(http.StatusOK, snapshot)
}
