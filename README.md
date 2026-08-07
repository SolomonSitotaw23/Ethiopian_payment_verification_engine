# Ethiopian Payment Receipt Verifier (Golang)

[![CI Workflow](https://github.com/abrhamyalew/telebirr-payment-verifier/actions/workflows/ci.yml/badge.svg)](https://github.com/abrhamyalew/telebirr-payment-verifier/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/abrhamyalew/telebirr-payment-verifier)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An enterprise-grade, high-performance Go service to verify payment receipts from **Telebirr**, **CBE (Commercial Bank of Ethiopia)**, **BOA (Bank of Abyssinia)**, and **Amhara Bank** against expected transaction details.

Built with **Gin**, native Go goroutine worker pools, structured logging, rate limiting, and an extensible provider plugin architecture.

## Supported Banks & Digital Wallets

- **✓ Telebirr**
- **✓ CBE (Commercial Bank of Ethiopia)**
- **✓ BOA (Bank of Abyssinia)**
- **✓ Amhara Bank**

---

## Features

- **Extensible Provider Architecture (`providers/`)**: Modular `Provider` interface (`Name()`, `ParseID()`, `CanHandle()`, `Verify()`) allowing easy community contributions for new Ethiopian banks.
- **Parallel Batch Processing**: Concurrently verifies up to 10+ receipts using native Go goroutines and semaphore worker pools.
- **Asynchronous Webhooks**: Process large verification batches asynchronously by passing a `callbackUrl`.
- **Health Checks & Observability**:
  - `GET /health` for Kubernetes / Docker liveness & readiness probes.
  - `GET /metrics` for tracking uptime, request counts, valid vs. invalid counts.
- **Structured JSON Logging**: Powered by Go standard `log/slog` with auto-generated `X-Request-ID` tracing.
- **API Protection & Rate Limiting**:
  - In-memory token-bucket rate limiter per client IP.
  - Optional `API_KEY` authentication (`X-API-Key` header).
- **Proxy Support**: Configurable HTTP/HTTPS proxy support for Telebirr requests.

---

## Getting Started

### Clone & Setup

```bash
git clone https://github.com/abrhamyalew/telebirr-payment-verifier.git
cd telebirr-payment-verifier
cp .env.example .env
go mod download
```

### Environment Variables

Configure your expected payment details in `.env`:

```env
PORT=5000
API_KEY=your-optional-secret-key
PROXY=proxy-host:proxy-port:proxy-username:proxy-password

# Telebirr
TELEBIRR_EXPECTED_AMOUNT=100
TELEBIRR_EXPECTED_RECIPIENT_ACCOUNT=1000123456789
TELEBIRR_EXPECTED_RECIPIENT_NAME=Abrham Yalew
TELEBIRR_EXPECTED_PAYMENT_YEAR=2025
TELEBIRR_EXPECTED_PAYMENT_MONTH=12
TELEBIRR_EXPECTED_STATUS=Completed

# CBE
CBE_EXPECTED_AMOUNT=40
CBE_EXPECTED_RECIPIENT_ACCOUNT=1****1234
CBE_EXPECTED_RECIPIENT_NAME=ABRHAM YALEW
CBE_EXPECTED_PAYMENT_YEAR=2025
CBE_EXPECTED_PAYMENT_MONTH=12
```

### Development & Build Commands

Using `make`:

```bash
make run          # Run server locally
make test         # Run unit test suite
make build        # Build static production binary
make docker-build # Build Docker container
make docker-run   # Run Docker container
```

---

## API Documentation

### 1. Health & Observability

- **`GET /health`**: Returns system status and uptime.
- **`GET /metrics`**: Returns metrics counters (total requests, valid receipts, failed receipts).

---

### 2. Verify Single Receipt

**POST** `/api/verify`

**Request:**
```json
{
  "receipt": "FT253183LQF089873510",
  "defaultVerification": true
}
```

**Response (200 OK):**
```json
{
  "message": "The receipt 'FT253183LQF089873510' is a valid receipt."
}
```

---

### 3. Verify Batch Receipts

**POST** `/api/verify/batch`

**Synchronous Request:**
```json
{
  "receipt": ["CJP9OSP9WZ", "FT25284X11PS79328"],
  "defaultVerification": true
}
```

**Asynchronous Request (Webhook Callback):**
```json
{
  "receipt": ["CJP9OSP9WZ", "FT25284X11PS79328"],
  "defaultVerification": true,
  "callbackUrl": "https://your-api.com/webhooks/payment-verifier"
}
```

**Async Response (202 Accepted):**
```json
{
  "message": "Batch processing started asynchronously. Verification result will be sent to callbackUrl.",
  "total": 2,
  "callbackUrl": "https://your-api.com/webhooks/payment-verifier"
}
```

---

## Adding a New Provider

To add support for a new bank or wallet:
1. Create a new file in `providers/yourbank.go`.
2. Implement the `providers.Provider` interface:
   ```go
   type Provider interface {
       Name() string
       ParseID(input string) string
       CanHandle(input string) bool
       Verify(ctx context.Context, input string, flags models.VerificationFlags, opts services.ReceiptRequestOptions) (string, error)
   }
   ```
3. Register your provider in `providers/init.go`:
   ```go
   DefaultRegistry.Register(NewYourBankProvider())
   ```

---

## License

[MIT](LICENSE) © Abrham Yalew
