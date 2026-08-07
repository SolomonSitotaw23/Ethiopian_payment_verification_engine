package models

type ParsedReceiptData struct {
	Amount        float64 `json:"amount"`
	RecipientName string  `json:"recipientName"`
	AccountNumber string  `json:"accountNumber"`
	Date          string  `json:"date"`
	Status        string  `json:"status,omitempty"`
}

type VerificationCheckResults struct {
	AmountMatched        bool `json:"amountMatched"`
	RecipientNameMatched bool `json:"recipientNameMatched"`
	AccountNumberMatched bool `json:"accountNumberMatched"`
	DateMatched          bool `json:"dateMatched"`
	StatusMatched        bool `json:"statusMatched"`
}

type DetailedVerifyResponse struct {
	Status     string                   `json:"status"` // "valid" or "mismatch"
	ReceiptID  string                   `json:"receiptId"`
	Provider   string                   `json:"provider"`
	Message    string                   `json:"message"`
	Parsed     ParsedReceiptData        `json:"parsed"`
	Checks     VerificationCheckResults `json:"checks"`
	Mismatches []string                 `json:"mismatches"`
}

type ErrorResponse struct {
	Error   string                  `json:"error"`
	Details *DetailedVerifyResponse `json:"details,omitempty"`
}

type BatchVerifyFailedItem struct {
	ReceiptID string                  `json:"receiptId"`
	Error     string                  `json:"error"`
	Details   *DetailedVerifyResponse `json:"details,omitempty"`
}

type BatchVerifySummary struct {
	Total   int `json:"total"`
	Valid   int `json:"valid"`
	Invalid int `json:"invalid"`
}

type BatchVerifyResponse struct {
	Result  []DetailedVerifyResponse `json:"result"`
	Failed  []BatchVerifyFailedItem  `json:"failed"`
	Summary BatchVerifySummary       `json:"summary"`
}

type AsyncBatchResponse struct {
	Message     string `json:"message"`
	Total       int    `json:"total"`
	CallbackURL string `json:"callbackUrl"`
}
