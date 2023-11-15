package repository

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/labstack/echo/v4"
)

type AccountRepository interface {
	Insert(ctx echo.Context, account entity.Account) error
	BatchInsert(ctx echo.Context, accounts []*entity.Account) error
	GetByAccountNumber(ctx echo.Context, accountNumber string) (*entity.Account, error)
	GetAll(ctx echo.Context) ([]*entity.Account, error)
}
