package validators

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"payment_verifier/config"
	"payment_verifier/models"
	"payment_verifier/services"
	"payment_verifier/utils"
	"github.com/ledongthuc/pdf"
)

var (
	cbeAmountRegex    = regexp.MustCompile(`(?i)Transferred Amount\s*([\d.]+\s*ETB)`)
	cbeDateRegex      = regexp.MustCompile(`(?i)Payment Date\s*&\s*Time\s*(\d{2}\/\d{2}\/\d{4})`)
	cbeAccountRegex   = regexp.MustCompile(`(?i)Receiver[\s\S]*?Account\s*(1\*{4}\d{4})`)
	cbeRecipientRegex = regexp.MustCompile(`(?i)Receiver\s*([A-Z\s]+?)\s*Account`)
)

func extractPdfText(pdfBytes []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(pdfBytes), int64(len(pdfBytes)))
	if err != nil {
		return "", err
	}
	plainTextReader, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}
	textBytes, err := io.ReadAll(plainTextReader)
	if err != nil {
		return "", err
	}
	return string(textBytes), nil
}

func extractRegexField(text string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func VerifyCBEDetailed(receiptID string, res *services.ReceiptDataResult, reqExpected *models.ExpectedDataRequest, flags models.VerificationFlags) (*models.DetailedVerifyResponse, error) {
	var rawAmt, dateStr, accountNo, recipientName string

	if res.CbeMB != nil {
		mb := res.CbeMB
		if len(mb.DateTimes) > 0 {
			datePart := strings.Split(mb.DateTimes[0], "T")[0]
			dParts := strings.Split(datePart, "-")
			if len(dParts) == 3 {
				year, month, day := dParts[0], dParts[1], dParts[2]
				dateStr = fmt.Sprintf("%s/%s/%s", month, day, year)
			}
		}
		rawAmt = mb.DebitAmount
		if rawAmt != "" {
			rawAmt = strings.Split(rawAmt, ".")[0]
		}
		accountNo = mb.CreditAccountNo
		recipientName = mb.CreditAccountHolder
	} else if len(res.CbePDFBytes) > 0 {
		text, err := extractPdfText(res.CbePDFBytes)
		if err != nil {
			return nil, utils.NewValidationError("Failed to parse CBE PDF receipt text")
		}
		pdfAmt := extractRegexField(text, cbeAmountRegex)
		if pdfAmt != "" {
			rawAmt = strings.Split(pdfAmt, ".")[0]
		}
		dateStr = extractRegexField(text, cbeDateRegex)
		accountNo = extractRegexField(text, cbeAccountRegex)
		recipientName = extractRegexField(text, cbeRecipientRegex)
	} else {
		return nil, utils.NewValidationError("No valid CBE receipt data found")
	}

	parsedAmt, _ := strconv.ParseFloat(rawAmt, 64)

	parsedData := models.ParsedReceiptData{
		Amount:        parsedAmt,
		RecipientName: recipientName,
		AccountNumber: accountNo,
		Date:          dateStr,
	}

	expected := utils.MergeExpectedData(reqExpected, config.Config.CBE.Expected)

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
		if parsedData.Date == "" {
			checks.DateMatched = false
			mismatches = append(mismatches, "No parsed data for date")
		} else {
			dParts := strings.Split(parsedData.Date, "/")
			if len(dParts) == 3 {
				month, _, year := dParts[0], dParts[1], dParts[2]
				if expected.PaymentYear != "" && year != expected.PaymentYear {
					checks.DateMatched = false
					mismatches = append(mismatches, fmt.Sprintf("Year mismatch. Expected: %s, Actual: %s", expected.PaymentYear, year))
				}
				if expected.PaymentMonth != "" && month != expected.PaymentMonth {
					checks.DateMatched = false
					mismatches = append(mismatches, fmt.Sprintf("Month mismatch. Expected: %s, Actual: %s", expected.PaymentMonth, month))
				}
			}
		}
	}

	// 2. Amount
	if shouldCheck(flags.Amount) {
		matched, errMsg := utils.CompareAmountFlexible(expected, rawAmt)
		if !matched {
			checks.AmountMatched = false
			mismatches = append(mismatches, errMsg)
		}
	}

	// 3. Recipient Name
	if shouldCheck(flags.RecipientName) {
		if expected.RecipientName != "" && strings.ToLower(strings.TrimSpace(expected.RecipientName)) != strings.ToLower(strings.TrimSpace(parsedData.RecipientName)) {
			checks.RecipientNameMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Mismatch on recipientName. Expected: %s, Actual: %s", expected.RecipientName, parsedData.RecipientName))
		}
	}

	// 4. Account Number
	if shouldCheck(flags.AccountNumber) {
		if expected.RecipientAccount != "" && strings.ToLower(strings.TrimSpace(expected.RecipientAccount)) != strings.ToLower(strings.TrimSpace(parsedData.AccountNumber)) {
			checks.AccountNumberMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Mismatch on accountNumber. Expected: %s, Actual: %s", expected.RecipientAccount, parsedData.AccountNumber))
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
		Provider:   "CBE",
		Message:    msg,
		Parsed:     parsedData,
		Checks:     checks,
		Mismatches: mismatches,
	}, nil
}
