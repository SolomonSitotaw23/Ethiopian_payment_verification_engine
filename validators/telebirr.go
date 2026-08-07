package validators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"payment_verifier/config"
	"payment_verifier/models"
	"payment_verifier/utils"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)

func normalizeStr(str string) string {
	return strings.ToLower(strings.TrimSpace(whitespaceRegex.ReplaceAllString(str, " ")))
}

func compareAmount(expected, parsed string) bool {
	expNum, expErr := strconv.ParseFloat(strings.TrimSpace(expected), 64)
	prsNum, prsErr := strconv.ParseFloat(strings.TrimSpace(parsed), 64)
	if expErr != nil || prsErr != nil {
		return strings.TrimSpace(expected) == strings.TrimSpace(parsed)
	}
	return expNum == prsNum
}

func VerifyTelebirr(rawHTML string, flags models.VerificationFlags) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return false, utils.NewValidationError("Invalid HTML content")
	}

	divText := strings.TrimSpace(doc.Find("div").Text())
	if strings.Contains(divText, "This request is not correct") {
		return false, utils.NewNotFoundError("Receipt not found or invalid")
	}

	accountAndName := doc.Find("#paid_reference_number").Text()
	parts := strings.Fields(accountAndName)
	var accountNumber, name string
	if len(parts) > 0 {
		accountNumber = parts[0]
		if len(parts) > 1 {
			name = strings.Join(parts[1:], " ")
		}
	}

	var invoiceTable *goquery.Selection
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(j int, td *goquery.Selection) {
			if strings.Contains(normalizeStr(td.Text()), "settled amount") {
				invoiceTable = s
			}
		})
	})

	var statusTable *goquery.Selection
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(j int, td *goquery.Selection) {
			if strings.Contains(normalizeStr(td.Text()), "transaction status") {
				statusTable = s
			}
		})
	})

	findColumnValueFromHeader := func(table *goquery.Selection, labelText string) string {
		if table == nil {
			return ""
		}
		matcher := normalizeStr(labelText)
		var headerTd *goquery.Selection
		table.Find("td").Each(func(i int, td *goquery.Selection) {
			if headerTd == nil && strings.Contains(normalizeStr(td.Text()), matcher) {
				headerTd = td
			}
		})
		if headerTd == nil {
			return ""
		}
		headerRow := headerTd.Parent()
		colIdx := -1
		headerRow.Find("td").Each(func(i int, td *goquery.Selection) {
			if colIdx == -1 && len(headerTd.Nodes) > 0 && len(td.Nodes) > 0 && td.Nodes[0] == headerTd.Nodes[0] {
				colIdx = i
			}
		})
		if colIdx == -1 {
			return ""
		}
		valueRow := headerRow.Next()
		valTd := valueRow.Find("td").Eq(colIdx)
		return strings.TrimSpace(whitespaceRegex.ReplaceAllString(valTd.Text(), " "))
	}

	findAdjacentValue := func(table *goquery.Selection, labelText string) string {
		if table == nil {
			return ""
		}
		matcher := normalizeStr(labelText)
		var targetTd *goquery.Selection
		table.Find("td").Each(func(i int, td *goquery.Selection) {
			if targetTd == nil && strings.Contains(normalizeStr(td.Text()), matcher) {
				targetTd = td
			}
		})
		if targetTd == nil {
			return ""
		}
		nextTd := targetTd.NextFiltered("td")
		if nextTd.Length() == 0 {
			nextTd = targetTd.Next()
		}
		return strings.TrimSpace(whitespaceRegex.ReplaceAllString(nextTd.Text(), " "))
	}

	amountRaw := ""
	date := ""
	if invoiceTable != nil {
		amountRaw = findColumnValueFromHeader(invoiceTable, "Settled Amount")
		date = findColumnValueFromHeader(invoiceTable, "Payment date")
	}

	amountFromTable := strings.TrimSpace(regexp.MustCompile(`(?i)Birr`).ReplaceAllString(amountRaw, ""))

	status := ""
	if statusTable != nil {
		status = findAdjacentValue(statusTable, "transaction status")
	}

	parsedData := models.TelebirrParsedData{
		Amount:        amountFromTable,
		Status:        status,
		RecipientName: name,
		Date:          date,
		AccountNumber: accountNumber,
	}

	expected := config.Config.Telebirr.Expected

	shouldCheck := func(fieldFlag *bool) bool {
		if flags.IsDefault {
			return true
		}
		return fieldFlag != nil && *fieldFlag
	}

	// 1. Amount
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

	// 2. Status
	if shouldCheck(flags.Status) {
		if expected.Status == "" {
			return false, utils.NewValidationError("No expected data for \"status\", failing verification.")
		}
		if parsedData.Status == "" {
			return false, utils.NewValidationError("No parsed data for \"status\", failing verification.")
		}
		if strings.TrimSpace(expected.Status) != strings.TrimSpace(parsedData.Status) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on status. Expected: %s, Actual: %s", expected.Status, parsedData.Status))
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
		if strings.TrimSpace(expected.RecipientName) != strings.TrimSpace(parsedData.RecipientName) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on recipientName. Expected: %s, Actual: %s", expected.RecipientName, parsedData.RecipientName))
		}
	}

	// 4. Account Number
	if shouldCheck(flags.AccountNumber) {
		if expected.AccountNumber == "" {
			return false, utils.NewValidationError("No expected data for \"accountNumber\", failing verification.")
		}
		if parsedData.AccountNumber == "" {
			return false, utils.NewValidationError("No parsed data for \"accountNumber\", failing verification.")
		}
		if strings.TrimSpace(expected.AccountNumber) != strings.TrimSpace(parsedData.AccountNumber) {
			return false, utils.NewValidationError(fmt.Sprintf("Mismatch on accountNumber. Expected: %s, Actual: %s", expected.AccountNumber, parsedData.AccountNumber))
		}
	}

	// 5. Date
	if shouldCheck(flags.Date) {
		if parsedData.Date == "" {
			return false, utils.NewValidationError("No parsed data for date")
		}
		dateParts := strings.Split(parsedData.Date, " ")
		if len(dateParts) > 0 {
			dParts := strings.Split(dateParts[0], "-")
			if len(dParts) >= 3 {
				_, month, year := dParts[0], dParts[1], dParts[2]
				if expected.PaymentYear != "" && year != expected.PaymentYear {
					return false, utils.NewValidationError(fmt.Sprintf("Year mismatch. Expected: %s, Actual: %s", expected.PaymentYear, year))
				}
				if expected.PaymentMonth != "" && month != expected.PaymentMonth {
					return false, utils.NewValidationError(fmt.Sprintf("Month mismatch. Expected: %s, Actual: %s", expected.PaymentMonth, month))
				}
			}
		}
	}

	return true, nil
}
