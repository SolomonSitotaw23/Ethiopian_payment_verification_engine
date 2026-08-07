package validators

import (
	"fmt"
	"strconv"
	"strings"

	"payment_verifier/config"
	"payment_verifier/models"
	"payment_verifier/utils"
)

func VerifyBOADetailed(receiptID string, parsedData *models.BoaParsedData, reqExpected *models.ExpectedDataRequest, flags models.VerificationFlags) (*models.DetailedVerifyResponse, error) {
	if parsedData == nil || parsedData.TransactionDate == "" {
		return nil, utils.NewValidationError("No parsed data for Transaction Date")
	}

	dateParts := strings.Split(parsedData.TransactionDate, " ")
	dParts := strings.Split(dateParts[0], "/")
	if len(dParts) < 3 {
		return nil, utils.NewValidationError("Invalid Transaction Date format")
	}

	_, month, year := dParts[0], dParts[1], dParts[2]
	amountStr := fmt.Sprintf("%v", parsedData.TransferredAmount)
	parsedAmt, _ := strconv.ParseFloat(amountStr, 64)

	parsedReceipt := models.ParsedReceiptData{
		Amount:        parsedAmt,
		RecipientName: parsedData.ReceiversName,
		AccountNumber: parsedData.ReceiversAccount,
		Date:          parsedData.TransactionDate,
	}

	expected := utils.MergeExpectedData(reqExpected, config.Config.BOA.Expected)

	checks := models.VerificationCheckResults{
		AmountMatched:        true,
		RecipientNameMatched: true,
		AccountNumberMatched: true,
		DateMatched:          true,
		StatusMatched:        true,
	}
	mismatches := make([]string, 0)

	shouldCheck := func(fieldFlag *bool) bool {
		if flags.IsDefault {
			return true
		}
		return fieldFlag != nil && *fieldFlag
	}

	// 1. Date
	if shouldCheck(flags.Date) {
		if expected.PaymentYear != "" && year != expected.PaymentYear {
			checks.DateMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Year mismatch. Expected: %s, Actual: %s", expected.PaymentYear, year))
		}
		if expected.PaymentMonth != "" && month != expected.PaymentMonth {
			checks.DateMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Month mismatch. Expected: %s, Actual: %s", expected.PaymentMonth, month))
		}
	}

	// 2. Amount
	if shouldCheck(flags.Amount) {
		matched, errMsg := utils.CompareAmountFlexible(expected, amountStr)
		if !matched {
			checks.AmountMatched = false
			mismatches = append(mismatches, errMsg)
		}
	}

	// 3. Recipient Name
	if shouldCheck(flags.RecipientName) {
		if expected.RecipientName != "" && strings.TrimSpace(expected.RecipientName) != strings.TrimSpace(parsedReceipt.RecipientName) {
			checks.RecipientNameMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Mismatch on recipientName. Expected: %s, Actual: %s", expected.RecipientName, parsedReceipt.RecipientName))
		}
	}

	// 4. Account Number
	if shouldCheck(flags.AccountNumber) {
		if expected.RecipientAccount != "" && strings.TrimSpace(expected.RecipientAccount) != strings.TrimSpace(parsedReceipt.AccountNumber) {
			checks.AccountNumberMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Mismatch on accountNumber. Expected: %s, Actual: %s", expected.RecipientAccount, parsedReceipt.AccountNumber))
		}
	}

	resStatus := "valid"
	msg := fmt.Sprintf("The receipt '%s' is a valid receipt.", receiptID)
	if len(mismatches) > 0 {
		resStatus = "mismatch"
		msg = mismatches[0]
	}

	return &models.DetailedVerifyResponse{
		Status:     resStatus,
		ReceiptID:  receiptID,
		Provider:   "BOA",
		Message:    msg,
		Parsed:     parsedReceipt,
		Checks:     checks,
		Mismatches: mismatches,
	}, nil
}
