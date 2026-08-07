package providers

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"payment_verifier/models"
	"payment_verifier/services"
	"payment_verifier/utils"
)

type Provider interface {
	Name() string
	ParseID(input string) string
	CanHandle(input string) bool
	Verify(ctx context.Context, input string, reqExpected *models.ExpectedDataRequest, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (*models.DetailedVerifyResponse, error)
}

type ProviderRegistry struct {
	mu        sync.RWMutex
	providers []Provider
}

var DefaultRegistry = NewRegistry()

func NewRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make([]Provider, 0),
	}
}

func (r *ProviderRegistry) Register(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = append(r.providers, p)
}

func (r *ProviderRegistry) FindProvider(input string, explicitProvider string) (Provider, string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	trimmed := strings.TrimSpace(input)
	exp := strings.ToLower(strings.TrimSpace(explicitProvider))

	if exp != "" && exp != "auto" {
		for _, p := range r.providers {
			pName := strings.ToLower(p.Name())
			if pName == exp || strings.Contains(pName, exp) {
				id := p.ParseID(trimmed)
				if id == "" {
					id = trimmed
				}
				return p, id
			}
		}
	}

	for _, p := range r.providers {
		if p.CanHandle(trimmed) {
			id := p.ParseID(trimmed)
			if id != "" {
				return p, id
			}
		}
	}
	return nil, ""
}

func (r *ProviderRegistry) VerifyReceipt(ctx context.Context, input string, explicitProvider string, reqExpected *models.ExpectedDataRequest, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (*models.DetailedVerifyResponse, error) {
	provider, id := r.FindProvider(input, explicitProvider)
	if provider == nil || id == "" {
		return nil, utils.NewValidationError(fmt.Sprintf("receipt '%s' is NOT a valid receipt", input))
	}
	return provider.Verify(ctx, input, reqExpected, flags, opts)
}
