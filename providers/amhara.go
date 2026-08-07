package providers

import (
	"context"
	"strings"

	"payment_verifier/models"
	"payment_verifier/services"
	"payment_verifier/utils"
	"payment_verifier/validators"
)

type AmharaProvider struct{}

func NewAmharaProvider() *AmharaProvider {
	return &AmharaProvider{}
}

func (p *AmharaProvider) Name() string {
	return "AmharaBank"
}

func (p *AmharaProvider) ParseID(input string) string {
	return utils.AmharaBankParser(input)
}

func (p *AmharaProvider) CanHandle(input string) bool {
	lower := strings.ToLower(strings.TrimSpace(input))
	if strings.Contains(lower, "amharabank") {
		return true
	}
	id := utils.AmharaBankParser(input)
	return id != ""
}

func (p *AmharaProvider) Verify(ctx context.Context, input string, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (string, error) {
	id := p.ParseID(input)
	if id == "" {
		return "", utils.NewAppError("Invalid Amhara Bannk Receipt ID", 400)
	}

	rawRes, err := services.GetReceiptData(id, services.ReceiptRequestOptions{})
	if err != nil {
		return "", err
	}
	if rawRes == nil || rawRes.AmharaData == nil {
		return "", utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
	}

	valid, err := validators.VerifyAmharaBank(rawRes.AmharaData, flags)
	if err != nil {
		return "", err
	}
	if valid {
		return id, nil
	}
	return "", utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
}
