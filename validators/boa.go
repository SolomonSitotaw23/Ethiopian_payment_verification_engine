package validators

import (
	"fmt"
	"strings"

	"payment_verifier/config"
	"payment_verifier/models"
	"payment_verifier/utils"
)

func VerifyBOA(parsedData *models.BoaParsedData, flags models.VerificationFlags) (bool, error) {
	if parsedData == nil || parsedData.TransactionDate == "" {
		return false, utils.NewValidationError("No parsed data for Transaction Date")
	}

	dateParts := strings.Split(parsedData.TransactionDate, " ")
	dParts := strings.Split(dateParts[0], "/")
	if len(dParts) < 3 {
		return false, utils.NewValidationError("Invalid Transaction Date format")
	}

	_, month, year := dParts[0], dParts[1], dParts[2]
	amountStr := fmt.Sprintf("%v", parsedData.TransferredAmount)

	receiptData := map[string]string{
		"amount":        amountStr,
		"month":         month,
		"year":          year,
		"recipientName": parsedData.ReceiversName,
		"accountNumber": parsedData.ReceiversAccount,
	}

	expected := config.Config.BOA.Expected

	shouldCheck := func(fieldFlag *bool) bool {
		if flags.IsDefault {
			return true
		}
		return fieldFlag != nil && *fieldFlag
	}

	// 1. Date
	if shouldCheck(flags.Date) {
		if expected.PaymentYear != "" && receiptData["year"] != expected.PaymentYear {
			return false, utils.NewValidationError(fmt.Sprintf("Year mismatch. Expected: %s, Actual: %s", expected.PaymentYear, receiptData["year"]))
		}
		if expected.PaymentMonth != "" && receiptData["month"] != expected.PaymentMonth {
			return false, utils.NewValidationError(fmt.Sprintf("Month mismatch. Expected: %s, Actual: %s", expected.PaymentMonth, receiptData["month"]))
		}
	}

	// 2. Amount
	if shouldCheck(flags.Amount) {
		if expected.Amount == "" {
			return false, utils.NewValidationError("No expected data for \"amount\", failing verification.")
		}
		if receiptData["amount"] == "" {
			return false, utils.NewValidationError("No parsed data for \"amount\", failing verification.")
		}
		if !compareAmount(expected.Amount, receiptData["amount"]) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on amount. Expected: %s, Actual: %s", expected.Amount, receiptData["amount"]))
		}
	}

	// 3. Recipient Name
	if shouldCheck(flags.RecipientName) {
		if expected.RecipientName == "" {
			return false, utils.NewValidationError("No expected data for \"recipientName\", failing verification.")
		}
		if receiptData["recipientName"] == "" {
			return false, utils.NewValidationError("No parsed data for \"recipientName\", failing verification.")
		}
		if strings.TrimSpace(expected.RecipientName) != strings.TrimSpace(receiptData["recipientName"]) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on recipientName. Expected: %s, Actual: %s", expected.RecipientName, receiptData["recipientName"]))
		}
	}

	// 4. Account Number
	if shouldCheck(flags.AccountNumber) {
		if expected.AccountNumber == "" {
			return false, utils.NewValidationError("No expected data for \"accountNumber\", failing verification.")
		}
		if receiptData["accountNumber"] == "" {
			return false, utils.NewValidationError("No parsed data for \"accountNumber\", failing verification.")
		}
		if strings.TrimSpace(expected.AccountNumber) != strings.TrimSpace(receiptData["accountNumber"]) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on accountNumber. Expected: %s, Actual: %s", expected.AccountNumber, receiptData["accountNumber"]))
		}
	}

	return true, nil
}
