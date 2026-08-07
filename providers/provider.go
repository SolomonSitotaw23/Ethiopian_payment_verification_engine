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
	Verify(ctx context.Context, input string, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (string, error)
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

func (r *ProviderRegistry) FindProvider(input string) (Provider, string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	trimmed := strings.TrimSpace(input)
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

func (r *ProviderRegistry) VerifyReceipt(ctx context.Context, input string, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (string, error) {
	provider, id := r.FindProvider(input)
	if provider == nil || id == "" {
		return "", utils.NewValidationError(fmt.Sprintf("receipt '%s' is NOT a valid receipt", input))
	}
	return provider.Verify(ctx, input, flags, opts)
}
