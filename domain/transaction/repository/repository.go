package repository

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Add(ctx echo.Context, transaction *entity.Transaction, trx *gorm.DB) *echo.HTTPError
	GetLastTransaction(ctx echo.Context, accountNumber string, transactionType *string, numberOfTransaction int) ([]*entity.Transaction, *echo.HTTPError)
}
