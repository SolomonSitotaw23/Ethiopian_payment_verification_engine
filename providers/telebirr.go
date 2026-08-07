package providers

import (
	"context"
	"strings"

	"payment_verifier/models"
	"payment_verifier/services"
	"payment_verifier/utils"
	"payment_verifier/validators"
)

type TelebirrProvider struct{}

func NewTelebirrProvider() *TelebirrProvider {
	return &TelebirrProvider{}
}

func (p *TelebirrProvider) Name() string {
	return "Telebirr"
}

func (p *TelebirrProvider) ParseID(input string) string {
	return utils.TelebirrParser(input)
}

func (p *TelebirrProvider) CanHandle(input string) bool {
	lower := strings.ToLower(strings.TrimSpace(input))
	if strings.Contains(lower, "ethiotelecom") {
		return true
	}
	id := utils.TelebirrParser(input)
	return id != ""
}

func (p *TelebirrProvider) Verify(ctx context.Context, input string, reqExpected *models.ExpectedDataRequest, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (*models.DetailedVerifyResponse, error) {
	id := p.ParseID(input)
	if id == "" {
		return nil, utils.NewAppError("Invalid TeleBirr Receipt ID", 400)
	}

	rawRes, err := services.GetReceiptData(id, opts)
	if err != nil {
		return nil, err
	}
	if rawRes == nil || rawRes.TelebirrHTML == nil {
		return nil, utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
	}

	return validators.VerifyTelebirrDetailed(id, *rawRes.TelebirrHTML, reqExpected, flags)
}
