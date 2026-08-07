package providers

import (
	"context"
	"strings"

	"payment_verifier/models"
	"payment_verifier/services"
	"payment_verifier/utils"
	"payment_verifier/validators"
)

type CBEProvider struct{}

func NewCBEProvider() *CBEProvider {
	return &CBEProvider{}
}

func (p *CBEProvider) Name() string {
	return "CBE"
}

func (p *CBEProvider) ParseID(input string) string {
	return utils.CBEParser(input)
}

func (p *CBEProvider) CanHandle(input string) bool {
	lower := strings.ToLower(strings.TrimSpace(input))
	if strings.Contains(lower, "cbe") {
		return true
	}
	id := utils.CBEParser(input)
	return id != ""
}

func (p *CBEProvider) Verify(ctx context.Context, input string, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (string, error) {
	id := p.ParseID(input)
	if id == "" {
		return "", utils.NewAppError("Invalid CBE Receipt ID", 400)
	}

	rawRes, err := services.GetReceiptData(id, services.ReceiptRequestOptions{})
	if err != nil {
		return "", err
	}
	if rawRes == nil || (rawRes.CbeMB == nil && len(rawRes.CbePDFBytes) == 0) {
		return "", utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
	}

	valid, err := validators.VerifyCBE(rawRes, flags)
	if err != nil {
		return "", err
	}
	if valid {
		return id, nil
	}
	return "", utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
}
