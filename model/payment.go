package model

type PaymentIntentRequest struct {
	BankBic          string                 `json:"bankBic"`
	BankAccount      string                 `json:"bankAccount"`
	Amount           float64                `json:"amount"`
	PaymentReference string                 `json:"paymentReference"`
	Description      string                 `json:"description"`
	FeeOnMerchant    bool                   `json:"feeOnMerchant"`
	ReturnUrl        string                 `json:"returnUrl"`
	CallbackUrl      string                 `json:"callbackUrl"`
	MetaData         map[string]interface{} `json:"meta_data"`
}

type PaymentIntentResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
