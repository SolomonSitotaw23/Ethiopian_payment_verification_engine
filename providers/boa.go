package providers

import (
	"context"
	"strings"

	"github.com/abrhamyalew/telebirr-payment-verifier/models"
	"github.com/abrhamyalew/telebirr-payment-verifier/services"
	"github.com/abrhamyalew/telebirr-payment-verifier/utils"
	"github.com/abrhamyalew/telebirr-payment-verifier/validators"
)

type BOAProvider struct{}

func NewBOAProvider() *BOAProvider {
	return &BOAProvider{}
}

func (p *BOAProvider) Name() string {
	return "BOA"
}

func (p *BOAProvider) ParseID(input string) string {
	return utils.BOAParser(input)
}

func (p *BOAProvider) CanHandle(input string) bool {
	lower := strings.ToLower(strings.TrimSpace(input))
	if strings.Contains(lower, "bankofabyssinia") {
		return true
	}
	id := utils.BOAParser(input)
	return id != ""
}

func (p *BOAProvider) Verify(ctx context.Context, input string, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (string, error) {
	id := p.ParseID(input)
	if id == "" {
		return "", utils.NewAppError("Invalid BOA Receipt ID", 400)
	}

	rawRes, err := services.GetReceiptData(id, services.ReceiptRequestOptions{})
	if err != nil {
		return "", err
	}
	if rawRes == nil || rawRes.BoaData == nil {
		return "", utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
	}

	valid, err := validators.VerifyBOA(rawRes.BoaData, flags)
	if err != nil {
		return "", err
	}
	if valid {
		return id, nil
	}
	return "", utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
}
