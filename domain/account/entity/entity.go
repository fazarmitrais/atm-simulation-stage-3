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
	Date          string  `json:"Date"`
}

func (a *Account) ToAccountResponse() *AccountResponse {
	return &AccountResponse{
		Name:          a.Name,
		AccountNumber: a.AccountNumber,
		Balance:       a.Balance,
		Date:          time.Now().Format("2006-01-02 15:04"),
	}
}
