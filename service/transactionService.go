package service

import (
	"time"

	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Service) CreateTransaction(c echo.Context, transaction trxEntity.Transaction, trx *gorm.DB) *echo.HTTPError {
	transaction.Date = time.Now()
	return s.TransactionRepository.Add(c, &transaction, trx)
}

func (s *Service) GetLastTransaction(c echo.Context, accountNumber string, transactionType *string, numOfLastTransaction int) ([]*trxEntity.Transaction, *echo.HTTPError) {
	trxs, err := s.TransactionRepository.GetLastTransaction(c, accountNumber, transactionType, numOfLastTransaction)
	if err != nil {
		return nil, err
	}
	return trxs, nil
}
