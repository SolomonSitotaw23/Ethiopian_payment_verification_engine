package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ExpectedData struct {
	Amount        string
	Status        string
	RecipientName string
	AccountNumber string
	PaymentYear   string
	PaymentMonth  string
}

type ProviderConfig struct {
	Expected ExpectedData
}

type AppConfig struct {
	CBE        ProviderConfig
	Telebirr   ProviderConfig
	BOA        ProviderConfig
	AmharaBank ProviderConfig
}

var Config AppConfig

func LoadConfig() {
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	Config = AppConfig{
		CBE: ProviderConfig{
			Expected: ExpectedData{
				Amount:        os.Getenv("CBE_EXPECTED_AMOUNT"),
				RecipientName: os.Getenv("CBE_EXPECTED_RECIPIENT_NAME"),
				AccountNumber: os.Getenv("CBE_EXPECTED_RECIPIENT_ACCOUNT"),
				PaymentYear:   os.Getenv("CBE_EXPECTED_PAYMENT_YEAR"),
				PaymentMonth:  os.Getenv("CBE_EXPECTED_PAYMENT_MONTH"),
			},
		},
		Telebirr: ProviderConfig{
			Expected: ExpectedData{
				Amount:        os.Getenv("TELEBIRR_EXPECTED_AMOUNT"),
				Status:        os.Getenv("TELEBIRR_EXPECTED_STATUS"),
				RecipientName: os.Getenv("TELEBIRR_EXPECTED_RECIPIENT_NAME"),
				AccountNumber: os.Getenv("TELEBIRR_EXPECTED_RECIPIENT_ACCOUNT"),
				PaymentYear:   os.Getenv("TELEBIRR_EXPECTED_PAYMENT_YEAR"),
				PaymentMonth:  os.Getenv("TELEBIRR_EXPECTED_PAYMENT_MONTH"),
			},
		},
		BOA: ProviderConfig{
			Expected: ExpectedData{
				Amount:        os.Getenv("BOA_EXPECTED_AMOUNT"),
				RecipientName: os.Getenv("BOA_EXPECTED_RECIPIENT_NAME"),
				AccountNumber: os.Getenv("BOA_EXPECTED_RECIPIENT_ACCOUNT"),
				PaymentYear:   os.Getenv("BOA_EXPECTED_PAYMENT_YEAR"),
				PaymentMonth:  os.Getenv("BOA_EXPECTED_PAYMENT_MONTH"),
			},
		},
		AmharaBank: ProviderConfig{
			Expected: ExpectedData{
				Amount:        os.Getenv("AB_EXPECTED_AMOUNT"),
				RecipientName: os.Getenv("AB_EXPECTED_RECIPIENT_NAME"),
				AccountNumber: os.Getenv("AB_EXPECTED_RECIPIENT_ACCOUNT"),
				PaymentYear:   os.Getenv("AB_EXPECTED_PAYMENT_YEAR"),
				PaymentMonth:  os.Getenv("AB_EXPECTED_PAYMENT_MONTH"),
			},
		},
	}
}
