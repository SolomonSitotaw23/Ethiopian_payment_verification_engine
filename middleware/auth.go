package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"payment_verifier/models"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		expectedKey := strings.TrimSpace(os.Getenv("API_KEY"))
		if expectedKey == "" {
			c.Next()
			return
		}

		clientKey := c.GetHeader("X-API-Key")
		if clientKey == "" {
			clientKey = c.Query("key")
		}

		if clientKey == "" || clientKey != expectedKey {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Unauthorized: Invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
