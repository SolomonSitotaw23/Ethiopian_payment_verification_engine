package validators

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/abrhamyalew/telebirr-payment-verifier/config"
	"github.com/abrhamyalew/telebirr-payment-verifier/models"
	"github.com/abrhamyalew/telebirr-payment-verifier/services"
	"github.com/abrhamyalew/telebirr-payment-verifier/utils"
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

func VerifyCBE(res *services.ReceiptDataResult, flags models.VerificationFlags) (bool, error) {
	var parsedData models.CbePdfParsedData

	if res.CbeMB != nil {
		mb := res.CbeMB
		var formattedDate string
		if len(mb.DateTimes) > 0 {
			datePart := strings.Split(mb.DateTimes[0], "T")[0]
			dParts := strings.Split(datePart, "-")
			if len(dParts) == 3 {
				year, month, day := dParts[0], dParts[1], dParts[2]
				formattedDate = fmt.Sprintf("%s/%s/%s", month, day, year)
			}
		}

		amt := mb.DebitAmount
		if amt != "" {
			amt = strings.Split(amt, ".")[0]
		}

		parsedData = models.CbePdfParsedData{
			Amount:        amt,
			Date:          formattedDate,
			AccountNumber: mb.CreditAccountNo,
			RecipientName: mb.CreditAccountHolder,
		}
	} else if len(res.CbePDFBytes) > 0 {
		text, err := extractPdfText(res.CbePDFBytes)
		if err != nil {
			return false, utils.NewValidationError("Failed to parse CBE PDF receipt text")
		}

		rawAmt := extractRegexField(text, cbeAmountRegex)
		amt := ""
		if rawAmt != "" {
			amt = strings.Split(rawAmt, ".")[0]
		}

		parsedData = models.CbePdfParsedData{
			Amount:        amt,
			Date:          extractRegexField(text, cbeDateRegex),
			AccountNumber: extractRegexField(text, cbeAccountRegex),
			RecipientName: extractRegexField(text, cbeRecipientRegex),
		}
	} else {
		return false, utils.NewValidationError("No valid CBE receipt data found")
	}

	expected := config.Config.CBE.Expected

	shouldCheck := func(fieldFlag *bool) bool {
		if flags.IsDefault {
			return true
		}
		return fieldFlag != nil && *fieldFlag
	}

	// 1. Date
	if shouldCheck(flags.Date) {
		if parsedData.Date == "" {
			return false, utils.NewValidationError("No parsed data for date")
		}
		dParts := strings.Split(parsedData.Date, "/")
		if len(dParts) == 3 {
			month, _, year := dParts[0], dParts[1], dParts[2]
			if expected.PaymentYear != "" && year != expected.PaymentYear {
				return false, utils.NewValidationError(fmt.Sprintf("Year mismatch. Expected: %s, Actual: %s", expected.PaymentYear, year))
			}
			if expected.PaymentMonth != "" && month != expected.PaymentMonth {
				return false, utils.NewValidationError(fmt.Sprintf("Month mismatch. Expected: %s, Actual: %s", expected.PaymentMonth, month))
			}
		}
	}

	// 2. Amount
	if shouldCheck(flags.Amount) {
		if expected.Amount == "" {
			return false, utils.NewValidationError("No expected data for \"amount\", failing verification.")
		}
		if parsedData.Amount == "" {
			return false, utils.NewValidationError("No parsed data for \"amount\", failing verification.")
		}
		if !compareAmount(expected.Amount, parsedData.Amount) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on amount. Expected: %s, Actual: %s", expected.Amount, parsedData.Amount))
		}
	}

	// 3. Recipient Name
	if shouldCheck(flags.RecipientName) {
		if expected.RecipientName == "" {
			return false, utils.NewValidationError("No expected data for \"recipientName\", failing verification.")
		}
		if parsedData.RecipientName == "" {
			return false, utils.NewValidationError("No parsed data for \"recipientName\", failing verification.")
		}
		if strings.ToLower(strings.TrimSpace(expected.RecipientName)) != strings.ToLower(strings.TrimSpace(parsedData.RecipientName)) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on recipientName. Expected: %s, Actual: %s", expected.RecipientName, parsedData.RecipientName))
		}
	}

	// 4. Account Number
	if shouldCheck(flags.AccountNumber) {
		if expected.AccountNumber == "" {
			return false, utils.NewValidationError("No expected data for \"accountNumber\", failing verification.")
		}
		if parsedData.AccountNumber == "" {
			return false, utils.NewValidationError("No expected data for \"accountNumber\", failing verification.")
		}
		if strings.ToLower(strings.TrimSpace(expected.AccountNumber)) != strings.ToLower(strings.TrimSpace(parsedData.AccountNumber)) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on accountNumber. Expected: %s, Actual: %s", expected.AccountNumber, parsedData.AccountNumber))
		}
	}

	return true, nil
}
