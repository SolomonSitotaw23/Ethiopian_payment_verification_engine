package services

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/abrhamyalew/telebirr-payment-verifier/utils"
)

var (
	defaultHTTPClient *http.Client
	telebirrProxyClient *http.Client
	proxyOnce         sync.Once
	proxyErr          error
)

func init() {
	defaultTransport := &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     60 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	defaultHTTPClient = &http.Client{
		Transport: defaultTransport,
		Timeout:   15 * time.Second,
	}
}

func parseProxyEnvToURL(raw string) (*url.URL, error) {
	trimmed := strings.TrimSpace(raw)
	parts := strings.Split(trimmed, ":")
	if len(parts) < 4 {
		return nil, utils.NewAppError("Proxy requested but PROXY env is invalid (expected host:port:user:pass)", 500)
	}

	host := parts[0]
	portStr := parts[1]
	username := parts[2]
	password := strings.Join(parts[3:], ":")

	if host == "" || portStr == "" || username == "" || password == "" {
		return nil, utils.NewAppError("Proxy requested but PROXY env is invalid (expected host:port:user:pass)", 500)
	}

	portNum, err := strconv.Atoi(portStr)
	if err != nil || portNum <= 0 || portNum > 65535 {
		return nil, utils.NewAppError("Proxy requested but PROXY env port is invalid", 500)
	}

	proxyURLStr := fmt.Sprintf("http://%s:%s@%s:%d", url.QueryEscape(username), url.QueryEscape(password), host, portNum)
	return url.Parse(proxyURLStr)
}

func getTelebirrProxyClient() (*http.Client, error) {
	proxyOnce.Do(func() {
		raw := os.Getenv("PROXY")
		if raw == "" {
			proxyErr = utils.NewAppError("Proxy requested but PROXY env is missing", 500)
			return
		}

		proxyURL, err := parseProxyEnvToURL(raw)
		if err != nil {
			proxyErr = err
			return
		}

		proxyTransport := &http.Transport{
			Proxy:               http.ProxyURL(proxyURL),
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
			IdleConnTimeout:     60 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		}

		telebirrProxyClient = &http.Client{
			Transport: proxyTransport,
			Timeout:   15 * time.Second,
		}
	})

	if proxyErr != nil {
		return nil, proxyErr
	}
	return telebirrProxyClient, nil
}
