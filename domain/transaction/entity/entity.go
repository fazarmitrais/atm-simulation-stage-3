package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID                      uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Type                    string    `json:"type"`
	AccountNumber           string    `json:"accountNumber"`
	TransferToAccountNumber string    `json:"transferToAccountNumber"`
	Amount                  float64   `json:"amount"`
	Date                    time.Time `json:"date"`
	ReferenceNumber         string    `json:"referenceNumber"`
}

const (
	TYPE_WITHDRAWAL = "WITHDRAWAL"
	TYPE_TRANSFER   = "TRANSFER"
)
