package models

import (
	"encoding/json"
	"fmt"
)

type VerificationFlags struct {
	IsDefault     bool  `json:"-"`
	Amount        *bool `json:"amount,omitempty"`
	Status        *bool `json:"status,omitempty"`
	RecipientName *bool `json:"recipientName,omitempty"`
	AccountNumber *bool `json:"accountNumber,omitempty"`
	Date          *bool `json:"date,omitempty"`
}

func (vf *VerificationFlags) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			vf.IsDefault = true
			return nil
		}
		return fmt.Errorf("defaultVerification cannot be false")
	}

	type Alias VerificationFlags
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*vf = VerificationFlags(aux)
	vf.IsDefault = false
	return nil
}

type ExpectedDataRequest struct {
	Amount           interface{} `json:"amount,omitempty"`
	MinAmount        interface{} `json:"minAmount,omitempty"`
	MaxAmount        interface{} `json:"maxAmount,omitempty"`
	RecipientAccount string      `json:"recipientAccount,omitempty"`
	RecipientName    string      `json:"recipientName,omitempty"`
	PaymentYear      string      `json:"paymentYear,omitempty"`
	PaymentMonth     string      `json:"paymentMonth,omitempty"`
	Status           string      `json:"status,omitempty"`
}

type VerifyRequestBody struct {
	Receipt             json.RawMessage      `json:"receipt" binding:"required"`
	Expected            *ExpectedDataRequest `json:"expected,omitempty"`
	DefaultVerification VerificationFlags    `json:"defaultVerification"`
	Strict              *bool                `json:"strict,omitempty"`
	Proxy               *bool                `json:"proxy,omitempty"`
	CallbackURL         string               `json:"callbackUrl,omitempty"`
}
