package model

type YayaCallbackRequest struct {
	PaymentId        string  `json:"paymentId"`
	PaymentReference string  `json:"paymentReference"`
	Amount           float64 `json:"amount"`
	Status           string  `json:"status"`
	TransactionId    string  `json:"transactionId"`
	BankBic          string  `json:"bankBic"`
	BankAccount      string  `json:"bankAccount"`
	PaymentCode      string  `json:"paymentCode"`
	Timestamp        string  `json:"timestamp"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}
