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
var regexpBirr = regexp.MustCompile(`(?i)Birr`)

func normalizeStr(str string) string {
	return strings.ToLower(strings.TrimSpace(whitespaceRegex.ReplaceAllString(str, " ")))
}

func VerifyTelebirrDetailed(receiptID, rawHTML string, reqExpected *models.ExpectedDataRequest, flags models.VerificationFlags) (*models.DetailedVerifyResponse, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return nil, utils.NewValidationError("Invalid HTML content")
	}

	divText := strings.TrimSpace(doc.Find("div").Text())
	if strings.Contains(divText, "This request is not correct") {
		return nil, utils.NewNotFoundError("Receipt not found or invalid")
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

	amountFromTable := strings.TrimSpace(regexpBirr.ReplaceAllString(amountRaw, ""))
	status := ""
	if statusTable != nil {
		status = findAdjacentValue(statusTable, "transaction status")
	}

	parsedAmt, _ := strconv.ParseFloat(amountFromTable, 64)

	parsedData := models.ParsedReceiptData{
		Amount:        parsedAmt,
		RecipientName: name,
		AccountNumber: accountNumber,
		Date:          date,
		Status:        status,
	}

	expected := utils.MergeExpectedData(reqExpected, config.Config.Telebirr.Expected)

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

	// 1. Amount
	if shouldCheck(flags.Amount) {
		matched, errMsg := utils.CompareAmountFlexible(expected, amountFromTable)
		if !matched {
			checks.AmountMatched = false
			mismatches = append(mismatches, errMsg)
		}
	}

	// 2. Status
	if shouldCheck(flags.Status) {
		if expected.Status != "" && strings.TrimSpace(expected.Status) != strings.TrimSpace(parsedData.Status) {
			checks.StatusMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Mismatch on status. Expected: %s, Actual: %s", expected.Status, parsedData.Status))
		}
	}

	// 3. Recipient Name
	if shouldCheck(flags.RecipientName) {
		if expected.RecipientName != "" && strings.TrimSpace(expected.RecipientName) != strings.TrimSpace(parsedData.RecipientName) {
			checks.RecipientNameMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Mismatch on recipientName. Expected: %s, Actual: %s", expected.RecipientName, parsedData.RecipientName))
		}
	}

	// 4. Account Number
	if shouldCheck(flags.AccountNumber) {
		if expected.RecipientAccount != "" && strings.TrimSpace(expected.RecipientAccount) != strings.TrimSpace(parsedData.AccountNumber) {
			checks.AccountNumberMatched = false
			mismatches = append(mismatches, fmt.Sprintf("Mismatch on accountNumber. Expected: %s, Actual: %s", expected.RecipientAccount, parsedData.AccountNumber))
		}
	}

	// 5. Date
	if shouldCheck(flags.Date) {
		if parsedData.Date == "" {
			checks.DateMatched = false
			mismatches = append(mismatches, "No parsed data for date")
		} else {
			dateParts := strings.Split(parsedData.Date, " ")
			if len(dateParts) > 0 {
				dParts := strings.Split(dateParts[0], "-")
				if len(dParts) >= 3 {
					_, month, year := dParts[0], dParts[1], dParts[2]
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
		Provider:   "Telebirr",
		Message:    msg,
		Parsed:     parsedData,
		Checks:     checks,
		Mismatches: mismatches,
	}, nil
}
