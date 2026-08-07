package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"payment_verifier/models"
	"payment_verifier/providers"
	"payment_verifier/services"
	"payment_verifier/utils"
)

func VerifyHandler(c *gin.Context) {
	var body models.VerifyRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var receiptStr string
	if err := json.Unmarshal(body.Receipt, &receiptStr); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "receipt must be a string"})
		return
	}

	useProxy := body.Proxy != nil && *body.Proxy
	isStrict := body.Strict == nil || *body.Strict // Default to strict true

	res, err := providers.DefaultRegistry.VerifyReceipt(c.Request.Context(), receiptStr, body.Expected, body.DefaultVerification, services.ReceiptRequestOptions{Proxy: useProxy})
	if err != nil {
		services.Metrics.RecordVerification(false)
		status := utils.GetHTTPStatus(err)
		c.JSON(status, models.ErrorResponse{Error: err.Error()})
		return
	}

	if res.Status == "valid" {
		services.Metrics.RecordVerification(true)
		c.JSON(http.StatusOK, res)
		return
	}

	// res.Status == "mismatch"
	services.Metrics.RecordVerification(false)
	if isStrict {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   res.Message,
			Details: res,
		})
		return
	}

	// Non-strict mode: 200 OK with mismatch details
	c.JSON(http.StatusOK, res)
}
