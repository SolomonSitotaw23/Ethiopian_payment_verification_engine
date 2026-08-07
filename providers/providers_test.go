package providers_test

import (
	"testing"

	"payment_verifier/providers"
)

func TestProviderRegistry(t *testing.T) {
	registry := providers.DefaultRegistry

	tests := []struct {
		name         string
		receiptInput string
		wantProvider string
		wantID       string
	}{
		{"Telebirr Plain", "CJP9OSP9WZ", "Telebirr", "CJP9OSP9WZ"},
		{"Telebirr URL", "https://transactioninfo.ethiotelecom.et/receipt/CJP9OSP9WZ", "Telebirr", "CJP9OSP9WZ"},
		{"CBE Query URL", "https://apps.cbe.com.et:100/?id=FT25292FRPWD89873710", "CBE", "FT25292FRPWD89873710"},
		{"CBE MB URL", "https://mbreciept.cbe.com.et/FT25292FRPWD-89873710", "CBE", "FT25292FRPWD-89873710"},
		{"BOA URL", "https://cs.bankofabyssinia.com/slip/?trx=FT25284X11PS79328", "BOA", "FT25284X11PS79328"},
		{"Amhara Bank Plain", "AB1234567890", "AmharaBank", "AB1234567890"},
		{"Unknown Invalid", "INVALID_RECEIPT_TEST_999", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, id := registry.FindProvider(tt.receiptInput)
			if tt.wantProvider == "" {
				if provider != nil {
					t.Errorf("Expected nil provider for %q, got %s", tt.receiptInput, provider.Name())
				}
			} else {
				if provider == nil {
					t.Fatalf("Expected provider %s for %q, got nil", tt.wantProvider, tt.receiptInput)
				}
				if provider.Name() != tt.wantProvider {
					t.Errorf("Provider.Name() = %s; want %s", provider.Name(), tt.wantProvider)
				}
				if id != tt.wantID {
					t.Errorf("Parsed ID = %s; want %s", id, tt.wantID)
				}
			}
		})
	}
}
