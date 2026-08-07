package utils_test

import (
	"testing"

	"payment_verifier/config"
	"payment_verifier/models"
	"payment_verifier/utils"
)

func TestMergeExpectedData(t *testing.T) {
	fallback := config.ExpectedData{
		Amount:        "100",
		AccountNumber: "0911223344",
		RecipientName: "Default Name",
	}

	// 1. Dynamic override test
	req := &models.ExpectedDataRequest{
		Amount:           230,
		RecipientAccount: "0988776655",
	}

	eff := utils.MergeExpectedData(req, fallback)

	if eff.Amount != 230 {
		t.Errorf("Expected dynamic amount 230, got %v", eff.Amount)
	}
	if eff.RecipientAccount != "0988776655" {
		t.Errorf("Expected dynamic account 0988776655, got %v", eff.RecipientAccount)
	}
	if eff.RecipientName != "Default Name" {
		t.Errorf("Expected fallback name 'Default Name', got %v", eff.RecipientName)
	}
}

func TestCompareAmountFlexible(t *testing.T) {
	// Exact match
	effExact := utils.EffectiveExpectedData{Amount: 230}
	matched, err := utils.CompareAmountFlexible(effExact, "230.00")
	if !matched {
		t.Errorf("Expected exact match for 230.00, got error: %s", err)
	}

	// Min amount match (success)
	effMin := utils.EffectiveExpectedData{MinAmount: 200}
	matched, err = utils.CompareAmountFlexible(effMin, "230.00")
	if !matched {
		t.Errorf("Expected min match for >=200 with 230.00, got error: %s", err)
	}

	// Min amount match (fail)
	matched, err = utils.CompareAmountFlexible(effMin, "150.00")
	if matched {
		t.Errorf("Expected min match failure for >=200 with 150.00")
	}

	// Range match (success)
	effRange := utils.EffectiveExpectedData{MinAmount: 200, MaxAmount: 300}
	matched, err = utils.CompareAmountFlexible(effRange, "250.00")
	if !matched {
		t.Errorf("Expected range match [200-300] for 250.00, got error: %s", err)
	}

	// Range match (fail)
	matched, err = utils.CompareAmountFlexible(effRange, "350.00")
	if matched {
		t.Errorf("Expected range match failure [200-300] for 350.00")
	}
}
