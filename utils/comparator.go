package utils

import (
	"fmt"
	"strconv"
	"strings"

	"payment_verifier/config"
	"payment_verifier/models"
)

type EffectiveExpectedData struct {
	Amount           interface{}
	MinAmount        interface{}
	MaxAmount        interface{}
	RecipientAccount string
	RecipientName    string
	PaymentYear      string
	PaymentMonth     string
	Status           string
}

func MergeExpectedData(req *models.ExpectedDataRequest, fallback config.ExpectedData) EffectiveExpectedData {
	eff := EffectiveExpectedData{
		Amount:           fallback.Amount,
		RecipientAccount: fallback.AccountNumber,
		RecipientName:    fallback.RecipientName,
		PaymentYear:      fallback.PaymentYear,
		PaymentMonth:     fallback.PaymentMonth,
		Status:           fallback.Status,
	}

	if req == nil {
		return eff
	}

	if req.Amount != nil {
		eff.Amount = req.Amount
	}
	if req.MinAmount != nil {
		eff.MinAmount = req.MinAmount
	}
	if req.MaxAmount != nil {
		eff.MaxAmount = req.MaxAmount
	}
	if req.RecipientAccount != "" {
		eff.RecipientAccount = req.RecipientAccount
	}
	if req.RecipientName != "" {
		eff.RecipientName = req.RecipientName
	}
	if req.PaymentYear != "" {
		eff.PaymentYear = req.PaymentYear
	}
	if req.PaymentMonth != "" {
		eff.PaymentMonth = req.PaymentMonth
	}
	if req.Status != "" {
		eff.Status = req.Status
	}

	return eff
}

func parseToFloat(val interface{}) (float64, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		str := strings.TrimSpace(v)
		if str == "" {
			return 0, false
		}
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, false
		}
		return num, true
	default:
		str := fmt.Sprintf("%v", v)
		num, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
		if err != nil {
			return 0, false
		}
		return num, true
	}
}

func CompareAmountFlexible(expected EffectiveExpectedData, parsedStr string) (bool, string) {
	actualNum, okActual := parseToFloat(parsedStr)
	if !okActual {
		return false, fmt.Sprintf("Parsed amount '%s' is not a valid number", parsedStr)
	}

	// 1. Min & Max Range check
	minNum, okMin := parseToFloat(expected.MinAmount)
	maxNum, okMax := parseToFloat(expected.MaxAmount)
	if okMin && okMax {
		if actualNum < minNum || actualNum > maxNum {
			return false, fmt.Sprintf("Amount out of expected range [%.2f - %.2f]. Actual: %.2f", minNum, maxNum, actualNum)
		}
		return true, ""
	}

	// 2. Min Amount check
	if okMin {
		if actualNum < minNum {
			return false, fmt.Sprintf("Amount below minimum expected %.2f. Actual: %.2f", minNum, actualNum)
		}
		return true, ""
	}

	// 3. Exact Amount check
	expNum, okExp := parseToFloat(expected.Amount)
	if okExp {
		if expNum != actualNum {
			return false, fmt.Sprintf("Mismatch on amount. Expected: %v, Actual: %.2f", expected.Amount, actualNum)
		}
		return true, ""
	}

	return false, "No valid expected amount criteria configured"
}
