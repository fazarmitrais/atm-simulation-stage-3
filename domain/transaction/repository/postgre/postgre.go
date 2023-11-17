package postgre

import (
	"errors"
	"net/http"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Postgre struct {
	db *gorm.DB
}

func NewPostgre(db *gorm.DB) *Postgre {
	return &Postgre{db}
}

func (p *Postgre) Add(ctx echo.Context, transaction *entity.Transaction, trx *gorm.DB) *echo.HTTPError {
	transaction.ID = uuid.New()
	res := trx.Create(transaction)
	if res.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error when inserting transaction : ", res.Error.Error())
	}
	return nil
}
func (p *Postgre) GetLastTransaction(ctx echo.Context, accountNumber string, transactionType *string, numberOfTransaction int) ([]*entity.Transaction, *echo.HTTPError) {
	var transactions []*entity.Transaction
	var res *gorm.DB
	if transactionType == nil {
		res = p.db.Model(&entity.Transaction{}).Where("account_number = ?", accountNumber).
			Order("Date DESC").
			Limit(numberOfTransaction).
			Find(&transactions)
	} else {
		res = p.db.Model(&entity.Transaction{}).Where("account_number = ? AND type = ?", accountNumber, *transactionType).
			Order("Date DESC").
			First(&transactions)
	}

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error get last transaction : "+res.Error.Error())
	}
	return transactions, nil
}
