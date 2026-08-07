package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"

	"github.com/abrhamyalew/telebirr-payment-verifier/models"
	"github.com/abrhamyalew/telebirr-payment-verifier/utils"
)

type ReceiptRequestOptions struct {
	Proxy bool
}

type ReceiptDataResult struct {
	TelebirrHTML *string
	CbeMB        *models.CbeMbParsedData
	CbePDFBytes  []byte
	BoaData      *models.BoaParsedData
	AmharaData   *models.AmharaBankParsedData
}

var (
	telebirrIDRegex = regexp.MustCompile(`^[A-Z0-9]{10}$`)
	cbeIDRegex      = regexp.MustCompile(`^[A-Z0-9]{12}\d{8}$`)
	cbeLegacyRegex   = regexp.MustCompile(`^[A-Z0-9]{12}&\d{8}$`)
	cbeMbRegex      = regexp.MustCompile(`^[A-Z0-9]{12}-\d{8}$`)
	boaIDRegex      = regexp.MustCompile(`^FT\d{5}[A-Z0-9]{5}\d{5}$`)
	amharaIDRegex   = regexp.MustCompile(`^[A-Z0-9]{12}$`)
)

func GetReceiptData(receiptID string, options ReceiptRequestOptions) (*ReceiptDataResult, error) {
	if telebirrIDRegex.MatchString(receiptID) {
		// Telebirr
		targetURL := fmt.Sprintf("https://transactioninfo.ethiotelecom.et/receipt/%s", receiptID)
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			return nil, utils.NewUpstreamServiceError("Failed to create Telebirr request")
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		client := defaultHTTPClient
		if options.Proxy {
			proxyClient, err := getTelebirrProxyClient()
			if err != nil {
				return nil, err
			}
			client = proxyClient
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, handleNetError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, utils.NewNotFoundError(fmt.Sprintf("Failed to fetch receipt. Status: %d", resp.StatusCode))
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, utils.NewUpstreamServiceError("Failed to read Telebirr response body")
		}

		htmlStr := string(bodyBytes)
		return &ReceiptDataResult{TelebirrHTML: &htmlStr}, nil

	} else if cbeIDRegex.MatchString(receiptID) || cbeLegacyRegex.MatchString(receiptID) || cbeMbRegex.MatchString(receiptID) {
		// CBE
		if cbeMbRegex.MatchString(receiptID) {
			targetURL := fmt.Sprintf("https://mb.cbe.com.et/api/v1/transactions/public/transaction-detail/%s", receiptID)
			req, err := http.NewRequest("GET", targetURL, nil)
			if err != nil {
				return nil, utils.NewUpstreamServiceError("Failed to create CBE MB request")
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("x-app-id", "d1292e42-7400-49de-a2d3-9731caa4c819")
			req.Header.Set("x-app-version", "0a01980b-9859-1369-8198-59f403820000")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

			resp, err := defaultHTTPClient.Do(req)
			if err != nil {
				return nil, handleNetError(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, utils.NewNotFoundError(fmt.Sprintf("Failed to fetch remote receipt. Status: %d", resp.StatusCode))
			}

			var mbData models.CbeMbParsedData
			if err := json.NewDecoder(resp.Body).Decode(&mbData); err != nil {
				return nil, utils.NewUpstreamServiceError("Failed to decode CBE MB JSON response")
			}
			return &ReceiptDataResult{CbeMB: &mbData}, nil
		}

		var targetURL string
		if cbeLegacyRegex.MatchString(receiptID) {
			targetURL = fmt.Sprintf("https://apps.cbe.com.et:100/BranchReceipt/%s", receiptID)
		} else {
			targetURL = fmt.Sprintf("https://apps.cbe.com.et:100/?id=%s", receiptID)
		}

		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			return nil, utils.NewUpstreamServiceError("Failed to create CBE request")
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		resp, err := defaultHTTPClient.Do(req)
		if err != nil {
			return nil, handleNetError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, utils.NewNotFoundError(fmt.Sprintf("Failed to fetch receipt. Status: %d", resp.StatusCode))
		}

		pdfBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, utils.NewUpstreamServiceError("Failed to read CBE PDF response")
		}

		return &ReceiptDataResult{CbePDFBytes: pdfBytes}, nil

	} else if boaIDRegex.MatchString(receiptID) {
		// BOA
		targetURL := fmt.Sprintf("https://cs.bankofabyssinia.com/api/onlineSlip/getDetails/?id=%s", receiptID)
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			return nil, utils.NewUpstreamServiceError("Failed to create BOA request")
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		resp, err := defaultHTTPClient.Do(req)
		if err != nil {
			return nil, handleNetError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, utils.NewNotFoundError(fmt.Sprintf("Failed to fetch receipt. Status: %d", resp.StatusCode))
		}

		var boaResp models.BoaApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&boaResp); err != nil {
			return nil, utils.NewUpstreamServiceError("Failed to decode BOA JSON response")
		}

		if len(boaResp.Body) == 0 {
			return nil, utils.NewNotFoundError("Receipt data not found in response")
		}

		return &ReceiptDataResult{BoaData: &boaResp.Body[0]}, nil

	} else if amharaIDRegex.MatchString(receiptID) {
		// Amhara Bank
		targetURL := fmt.Sprintf("https://transaction.amharabank.com.et/%s", receiptID)
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			return nil, utils.NewUpstreamServiceError("Failed to create Amhara Bank request")
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		resp, err := defaultHTTPClient.Do(req)
		if err != nil {
			return nil, handleNetError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, utils.NewNotFoundError("Receipt data not found or invalid")
		}

		var amharaResp models.AmharaBankApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&amharaResp); err != nil {
			return nil, utils.NewNotFoundError("Receipt data not found or invalid")
		}

		if !amharaResp.Status {
			return nil, utils.NewNotFoundError("Receipt data not found or invalid")
		}

		return &ReceiptDataResult{AmharaData: &amharaResp.Data}, nil
	}

	return nil, utils.NewValidationError(fmt.Sprintf("receipt '%s' is NOT a valid receipt", receiptID))
}

func handleNetError(err error) error {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return utils.NewConnectionTimeoutError("Upstream receipt service timed out")
	}
	return utils.NewUpstreamServiceError("Upstream receipt service is unavailable")
}
