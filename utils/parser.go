package utils

import "net/url"

import (
	"regexp"
	"strings"
)

var (
	telebirrPattern  = regexp.MustCompile(`^[A-Z0-9]{10}$`)
	cbeQueryPattern  = regexp.MustCompile(`^[A-Z0-9]{12}\d{8}$`)
	cbePathPattern   = regexp.MustCompile(`^[A-Z0-9]{12}&\d{8}$`)
	cbeMbPattern     = regexp.MustCompile(`^[A-Z0-9]{12}-\d{8}$`)
	boaPattern       = regexp.MustCompile(`^FT\d{5}[A-Z0-9]{5}\d{5}$`)
	amharaPattern    = regexp.MustCompile(`^[A-Z0-9]{12}$`)
)

func TelebirrParser(input string) string {
	trimInput := strings.TrimSpace(input)
	if trimInput == "" {
		return ""
	}

	id := trimInput
	if strings.HasPrefix(trimInput, "https://") {
		parts := strings.Split(trimInput, "/receipt/")
		if len(parts) > 1 {
			id = parts[1]
		}
	}

	if telebirrPattern.MatchString(id) {
		return id
	}
	return ""
}

func CBEParser(input string) string {
	trimInput := strings.TrimSpace(input)
	if trimInput == "" {
		return ""
	}

	parsedURL, err := url.Parse(trimInput)
	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		if parsedURL.Query().Has("id") {
			id := strings.TrimSpace(parsedURL.Query().Get("id"))
			if cbeQueryPattern.MatchString(id) {
				return id
			}
		} else {
			if strings.Contains(trimInput, "mbreciept.cbe.com.et") || strings.Contains(trimInput, "mb.cbe.com.et/api/v1/transactions/public/transaction-detail/") {
				parts := strings.Split(trimInput, "/")
				if len(parts) > 0 {
					id := parts[len(parts)-1]
					if cbeMbPattern.MatchString(id) {
						return id
					}
				}
			}

			if strings.HasPrefix(trimInput, "https") {
				parts := strings.Split(trimInput, "/BranchReceipt/")
				if len(parts) > 1 {
					id := parts[1]
					if cbePathPattern.MatchString(id) {
						return id
					}
				}
			}
		}
	}

	if cbeQueryPattern.MatchString(trimInput) || cbePathPattern.MatchString(trimInput) || cbeMbPattern.MatchString(trimInput) {
		return trimInput
	}

	return ""
}

func BOAParser(input string) string {
	trimInput := strings.TrimSpace(input)
	if trimInput == "" {
		return ""
	}

	parsedURL, err := url.Parse(trimInput)
	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		receiptNumber := parsedURL.Query().Get("trx")
		if receiptNumber != "" && boaPattern.MatchString(receiptNumber) {
			return receiptNumber
		}
	}

	if boaPattern.MatchString(trimInput) {
		return trimInput
	}

	return ""
}

func AmharaBankParser(input string) string {
	trimInput := strings.TrimSpace(input)
	if trimInput == "" {
		return ""
	}

	parts := strings.Split(trimInput, "/")
	if len(parts) > 0 {
		id := parts[len(parts)-1]
		if amharaPattern.MatchString(id) {
			return id
		}
	}

	if amharaPattern.MatchString(trimInput) {
		return trimInput
	}

	return ""
}
