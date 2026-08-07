package services

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"payment_verifier/models"
)

func SendWebhookCallback(callbackURL string, data models.BatchVerifyResponse) {
	go func() {
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			slog.Error("Failed to marshal webhook payload", "error", err, "url", callbackURL)
			return
		}

		req, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			slog.Error("Failed to create webhook request", "error", err, "url", callbackURL)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Ethiopian-Payment-Receipt-Verifier/1.0")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Failed to send webhook callback", "error", err, "url", callbackURL)
			return
		}
		defer resp.Body.Close()

		slog.Info("Webhook callback sent", "url", callbackURL, "status", resp.StatusCode)
	}()
}
