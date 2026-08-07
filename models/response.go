package models

type SingleVerifyResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type BatchVerifyFailedItem struct {
	ReceiptID string `json:"receiptId"`
	Error     string `json:"error"`
}

type BatchVerifySummary struct {
	Total   int `json:"total"`
	Valid   int `json:"valid"`
	Invalid int `json:"invalid"`
}

type BatchVerifyResponse struct {
	Result  []string                `json:"result"`
	Failed  []BatchVerifyFailedItem `json:"failed"`
	Summary BatchVerifySummary      `json:"summary"`
}
