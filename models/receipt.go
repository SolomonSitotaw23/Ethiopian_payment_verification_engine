package models

type TelebirrParsedData struct {
	Amount        string `json:"amount"`
	Status        string `json:"status"`
	RecipientName string `json:"recipientName"`
	Date          string `json:"date"`
	AccountNumber string `json:"accountNumber"`
}

type CbeMbParsedData struct {
	ID                  string   `json:"id"`
	DebitAmount         string   `json:"debitAmount"`
	AmountDebited       string   `json:"amountDebited"`
	CreditAccountNo     string   `json:"creditAccountNo"`
	CreditAccountHolder string   `json:"creditAccountHolder"`
	DateTimes           []string `json:"dateTimes"`
}

type CbePdfParsedData struct {
	Amount        string `json:"amount"`
	Date          string `json:"date"`
	AccountNumber string `json:"accountNumber"`
	RecipientName string `json:"recipientName"`
}

type BoaParsedData struct {
	TransactionDate  string      `json:"Transaction Date"`
	TransferredAmount interface{} `json:"Transferred Amount"`
	ReceiversName    string      `json:"Receiver's Name"`
	ReceiversAccount string      `json:"Receiver's Account"`
}

type BoaApiResponse struct {
	Body []BoaParsedData `json:"body"`
}

type AmharaBankParsedData struct {
	Status          string `json:"status"`
	BookingDate     string `json:"bookingDate"`
	Amount          string `json:"amount"`
	CreditorName    string `json:"creditorName"`
	CreditAccountId string `json:"creditAccountId"`
}

type AmharaBankApiResponse struct {
	Status bool                 `json:"status"`
	Data   AmharaBankParsedData `json:"data"`
}
