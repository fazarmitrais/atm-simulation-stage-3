package postgre

import (
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

func (p *Postgre) Insert(ctx echo.Context, account entity.Account) error {
	account.CreatedAt = time.Now()
	return p.db.Create(account).Error
}

func (p *Postgre) BatchInsert(ctx echo.Context, accounts []*entity.Account) error {
	now := time.Now()
	for _, a := range accounts {
		a.CreatedAt = now
	}
	return p.db.Create(&accounts).Error
}

func (p *Postgre) GetByAccountNumber(ctx echo.Context, accountNumber string) (*entity.Account, error) {
	var account entity.Account
	tx := p.db.Where("account_number = ?", accountNumber).Find(&account)
	if tx.Error == gorm.ErrRecordNotFound || tx.RowsAffected == 0 {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &account, nil
}

func (p *Postgre) GetAll(ctx echo.Context) ([]*entity.Account, error) {
	var accounts []*entity.Account
	tx := p.db.Order("created_at DESC").Find(&accounts)
	if tx.Error == gorm.ErrEmptySlice || tx.RowsAffected == 0 {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return accounts, nil
}
