package providers

import (
	"context"
	"strings"

	"payment_verifier/models"
	"payment_verifier/services"
	"payment_verifier/utils"
	"payment_verifier/validators"
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

func (p *BOAProvider) Verify(ctx context.Context, input string, reqExpected *models.ExpectedDataRequest, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (*models.DetailedVerifyResponse, error) {
	id := p.ParseID(input)
	if id == "" {
		return nil, utils.NewAppError("Invalid BOA Receipt ID", 400)
	}

	rawRes, err := services.GetReceiptData(id, services.ReceiptRequestOptions{})
	if err != nil {
		return nil, err
	}
	if rawRes == nil || rawRes.BoaData == nil {
		return nil, utils.NewValidationError("receipt '" + input + "' is NOT a valid receipt")
	}

	return validators.VerifyBOADetailed(id, rawRes.BoaData, reqExpected, flags)
}
