package entity

import "time"

type Account struct {
	Name          string    `json:"name" gorm:"primaryKey"`
	AccountNumber string    `json:"accountNumber"`
	PIN           string    `json:"pin"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"createdAt"`
}

type AccountResponse struct {
	Name          string  `json:"name"`
	AccountNumber string  `json:"accountNumber"`
	Balance       float64 `json:"balance"`
}

type Transfer struct {
	FromAccountNumber string  `json:"fromAccountNumber"`
	ToAccountNumber   string  `json:"toAccountNumber"`
	ReferenceNumber   string  `json:"referenceNumber"`
	Amount            float64 `json:"amount"`
}

func (a *Account) ToAccountResponse() *AccountResponse {
	return &AccountResponse{
		Name:          a.Name,
		AccountNumber: a.AccountNumber,
		Balance:       a.Balance,
	}
}
