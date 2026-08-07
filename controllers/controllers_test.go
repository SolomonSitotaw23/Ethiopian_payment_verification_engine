package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/abrhamyalew/telebirr-payment-verifier/controllers"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", controllers.HealthHandler)
	r.GET("/metrics", controllers.MetricsHandler)
	r.POST("/api/verify", controllers.VerifyHandler)
	r.POST("/api/verify/batch", controllers.BatchVerifyHandler)
	return r
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if res["status"] != "UP" {
		t.Errorf("Expected status UP, got %v", res["status"])
	}
}

func TestMetricsEndpoint(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestInvalidReceiptVerification(t *testing.T) {
	router := setupTestRouter()
	w := httptest.NewRecorder()
	payload := `{"receipt": "INVALID_RECEIPT_TEST_999", "defaultVerification": true}`
	req, _ := http.NewRequest("POST", "/api/verify", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid receipt format, got %d", w.Code)
	}
}
