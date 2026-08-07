package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/abrhamyalew/telebirr-payment-verifier/models"
	"github.com/abrhamyalew/telebirr-payment-verifier/providers"
	"github.com/abrhamyalew/telebirr-payment-verifier/services"
	"github.com/abrhamyalew/telebirr-payment-verifier/utils"
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
	validID, err := providers.DefaultRegistry.VerifyReceipt(c.Request.Context(), receiptStr, body.DefaultVerification, services.ReceiptRequestOptions{Proxy: useProxy})
	if err != nil {
		services.Metrics.RecordVerification(false)
		status := utils.GetHTTPStatus(err)
		c.JSON(status, models.ErrorResponse{Error: err.Error()})
		return
	}

	services.Metrics.RecordVerification(true)
	c.JSON(http.StatusOK, models.SingleVerifyResponse{
		Message: fmt.Sprintf("The receipt '%s' is a valid receipt.", validID),
	})
}
