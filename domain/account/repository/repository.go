package repository

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateTransaction() *gorm.DB
	Insert(ctx echo.Context, account entity.Account) *echo.HTTPError
	UpdateBalance(ctx echo.Context, account entity.Account, trx *gorm.DB) *echo.HTTPError
	BatchInsert(ctx echo.Context, accounts []*entity.Account) *echo.HTTPError
	GetByAccountNumber(ctx echo.Context, accountNumber string) (*entity.Account, *echo.HTTPError)
	GetAll(ctx echo.Context) ([]*entity.Account, *echo.HTTPError)
}
