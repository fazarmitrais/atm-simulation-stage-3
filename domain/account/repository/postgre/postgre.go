package postgre

import (
	"net/http"
	"time"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Postgre struct {
	db *gorm.DB
}

func NewPostgre(db *gorm.DB) *Postgre {
	return &Postgre{db}
}

func (p *Postgre) CreateTransaction() *gorm.DB {
	return p.db.Begin()
}

func (p *Postgre) Insert(ctx echo.Context, account entity.Account) *echo.HTTPError {
	account.CreatedAt = time.Now()
	err := p.db.Create(account).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting account : "+err.Error())
	}
	return nil
}

func (p *Postgre) UpdateBalance(ctx echo.Context, account entity.Account, trx *gorm.DB) *echo.HTTPError {
	res := trx.Model(&account).Update("balance", account.Balance)
	if res.Error != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "Error updating balance "+res.Error.Error())
	}
	return nil
}

func (p *Postgre) BatchInsert(ctx echo.Context, accounts []*entity.Account) *echo.HTTPError {
	now := time.Now()
	for _, a := range accounts {
		a.CreatedAt = now
	}
	err := p.db.Create(&accounts).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error batch insert account : "+err.Error())
	}
	return nil
}

func (p *Postgre) GetByAccountNumber(ctx echo.Context, accountNumber string) (*entity.Account, *echo.HTTPError) {
	var account entity.Account
	tx := p.db.Where("account_number = ?", accountNumber).Find(&account)
	if tx.Error == gorm.ErrRecordNotFound || tx.RowsAffected == 0 {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error getting account by account number : "+tx.Error.Error())
	}
	return &account, nil
}

func (p *Postgre) GetAll(ctx echo.Context) ([]*entity.Account, *echo.HTTPError) {
	var accounts []*entity.Account
	tx := p.db.Order("created_at DESC").Find(&accounts)
	if tx.Error == gorm.ErrEmptySlice || tx.RowsAffected == 0 {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error getting all account : "+tx.Error.Error())
	}
	return accounts, nil
}
