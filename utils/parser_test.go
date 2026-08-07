package utils_test

import (
	"testing"

	"github.com/abrhamyalew/telebirr-payment-verifier/utils"
)

func TestTelebirrParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Plain ID", "CJP9OSP9WZ", "CJP9OSP9WZ"},
		{"Full URL", "https://transactioninfo.ethiotelecom.et/receipt/CJP9OSP9WZ", "CJP9OSP9WZ"},
		{"Invalid ID", "INVALID12345", ""},
		{"Empty String", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.TelebirrParser(tt.input)
			if got != tt.expected {
				t.Errorf("TelebirrParser(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCBEParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Query Param URL", "https://apps.cbe.com.et:100/?id=FT25292FRPWD89873710", "FT25292FRPWD89873710"},
		{"Branch Receipt URL", "https://apps.cbe.com.et:100/BranchReceipt/FT25292FRPWD&89873710", "FT25292FRPWD&89873710"},
		{"MB Receipt URL", "https://mbreciept.cbe.com.et/FT25292FRPWD-89873710", "FT25292FRPWD-89873710"},
		{"Plain Query Pattern", "FT25292FRPWD89873710", "FT25292FRPWD89873710"},
		{"Plain Path Pattern", "FT25292FRPWD&89873710", "FT25292FRPWD&89873710"},
		{"Plain MB Pattern", "FT25292FRPWD-89873710", "FT25292FRPWD-89873710"},
		{"Invalid", "INVALID", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.CBEParser(tt.input)
			if got != tt.expected {
				t.Errorf("CBEParser(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBOAParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Full URL", "https://cs.bankofabyssinia.com/slip/?trx=FT25284X11PS79328", "FT25284X11PS79328"},
		{"Plain ID", "FT25284X11PS79328", "FT25284X11PS79328"},
		{"Invalid", "FT123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.BOAParser(tt.input)
			if got != tt.expected {
				t.Errorf("BOAParser(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestAmharaBankParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Full URL", "https://transaction.amharabank.com.et/AB1234567890", "AB1234567890"},
		{"Plain ID", "AB1234567890", "AB1234567890"},
		{"Invalid", "AB123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.AmharaBankParser(tt.input)
			if got != tt.expected {
				t.Errorf("AmharaBankParser(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
