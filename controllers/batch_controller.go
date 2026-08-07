package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"payment_verifier/config"
	"payment_verifier/models"
	"payment_verifier/providers"
	"payment_verifier/services"
)

type AsyncBatchResponse struct {
	Message     string `json:"message"`
	Total       int    `json:"total"`
	CallbackURL string `json:"callbackUrl"`
}

func BatchVerifyHandler(c *gin.Context) {
	var body models.VerifyRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var receipts []string
	if err := json.Unmarshal(body.Receipt, &receipts); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "receipt must be an array of strings"})
		return
	}

	useProxy := body.Proxy != nil && *body.Proxy
	ctx := c.Request.Context()

	processor := func(item string) (string, error) {
		validID, err := providers.DefaultRegistry.VerifyReceipt(ctx, item, body.DefaultVerification, services.ReceiptRequestOptions{Proxy: useProxy})
		if err != nil {
			services.Metrics.RecordVerification(false)
			return "", err
		}
		services.Metrics.RecordVerification(true)
		return validID, nil
	}

	callbackURL := strings.TrimSpace(body.CallbackURL)
	if callbackURL != "" {
		// Asynchronous mode
		go func() {
			results := services.ProcessBatch(receipts, processor, config.Performance.DefaultConcurrency)
			services.SendWebhookCallback(callbackURL, results)
		}()

		c.JSON(http.StatusAccepted, AsyncBatchResponse{
			Message:     "Batch processing started asynchronously. Verification result will be sent to callbackUrl.",
			Total:       len(receipts),
			CallbackURL: callbackURL,
		})
		return
	}

	// Synchronous mode
	results := services.ProcessBatch(receipts, processor, config.Performance.DefaultConcurrency)
	c.JSON(http.StatusOK, results)
}
