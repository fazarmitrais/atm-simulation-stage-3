package service

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	trxRepository "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/repository"
	"github.com/labstack/echo/v4"
)

type Service struct {
	AccountRepository     repository.AccountRepository
	AccountCsvRepository  repository.AccountCsvRepository
	TransactionRepository trxRepository.TransactionRepository
}

func NewService(accountRepository repository.AccountRepository, accountCsvRepository repository.AccountCsvRepository, transactionRepository trxRepository.TransactionRepository) *Service {
	return &Service{accountRepository, accountCsvRepository, transactionRepository}
}

type ServiceInterface interface {
	Insert(ctx echo.Context, account entity.Account) error
	PINValidation(c echo.Context, account entity.Account) *echo.HTTPError
	Transfer(ctx, transfer entity.Transfer) (*entity.Account, *echo.HTTPError)
	BalanceCheck(ctx echo.Context, acctNbr string) (*entity.Account, *echo.HTTPError)
	Import(path string) error
	GetAll() ([]*entity.Account, error)
	GetByAccountNumber(c echo.Context, acctNbr string) (*entity.Account, *echo.HTTPError)
	GetLastTransaction(c echo.Context, accountNumber string, transactionType *string, numOfLastTransaction int) ([]*trxEntity.Transaction, error)
}
