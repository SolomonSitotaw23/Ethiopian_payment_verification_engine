package providers

import (
	"context"
	"regexp"
	"strings"

	"github.com/abrhamyalew/telebirr-payment-verifier/models"
	"github.com/abrhamyalew/telebirr-payment-verifier/services"
	"github.com/abrhamyalew/telebirr-payment-verifier/utils"
	"github.com/abrhamyalew/telebirr-payment-verifier/validators"
)

var telebirrRegex = regexp.MustCompile(`^[A-Z0-9]{10}$`)

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

func (p *TelebirrProvider) Verify(ctx context.Context, input string, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (string, error) {
	id := p.ParseID(input)
	if id == "" {
		return "", utils.NewAppError("Invalid TeleBirr Receipt ID", 400)
	}

	rawRes, err := services.GetReceiptData(id, opts)
	if err != nil {
		return "", err
	}
	if rawRes == nil || rawRes.TelebirrHTML == nil {
		return "", utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
	}

	valid, err := validators.VerifyTelebirr(*rawRes.TelebirrHTML, flags)
	if err != nil {
		return "", err
	}
	if valid {
		return id, nil
	}
	return "", utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
}
