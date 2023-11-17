package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/labstack/echo/v4"
)

func (s *Service) PINValidation(c echo.Context, account entity.Account) *echo.HTTPError {
	if strings.Trim(account.AccountNumber, " ") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Account Number is required")
	} else if strings.Trim(account.PIN, " ") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "PIN is required")
	} else if len(account.AccountNumber) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Account Number should have 6 digits length")
	} else if len(account.PIN) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "PIN should have 6 digits length")
	} else if _, err := strconv.Atoi(account.AccountNumber); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Account Number should only contains numbers")
	} else if _, err := strconv.Atoi(account.PIN); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "PIN should only contains numbers")
	}
	accFromDb, err := s.AccountRepository.GetByAccountNumber(c, account.AccountNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get account")
	}
	if accFromDb == nil || accFromDb.PIN != account.PIN {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Account Number/PIN")
	}
	return nil
}

func (s *Service) Withdraw(ctx echo.Context, accountNumber string, withdrawAmount float64) (*entity.AccountResponse, *echo.HTTPError) {
	accFromDb, err := s.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	if withdrawAmount <= 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid withdraw amount")
	} else if withdrawAmount > 1000 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Maximum amount to withdraw is $1000")
	} else if int(withdrawAmount)%10 != 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid ammount")
	}

	if accFromDb.Balance < withdrawAmount {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Insufficient balance $%0.f", withdrawAmount))
	}
	accFromDb.Balance -= withdrawAmount
	trx := s.AccountRepository.CreateTransaction()
	err = s.AccountRepository.UpdateBalance(ctx, *accFromDb, trx)
	if err != nil {
		return nil, err
	}
	errl := s.CreateTransactionHistory(ctx, trxEntity.Transaction{
		AccountNumber: accountNumber,
		Amount:        withdrawAmount,
		Type:          trxEntity.TYPE_WITHDRAWAL,
	}, trx)
	trx.Commit()
	if err != nil {
		fmt.Printf("Error when creating transaction history: %v \n", errl)
	}
	return accFromDb.ToAccountResponse(), nil
}

func (s *Service) GetByAccountNumber(c echo.Context, acctNbr string) (*entity.Account, *echo.HTTPError) {
	if strings.Trim(acctNbr, " ") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Account Number is required")
	} else if len(acctNbr) < 6 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Account Number should have 6 digits length")
	} else if _, err := strconv.Atoi(acctNbr); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Account Number should only contains numbers")
	}
	accFromDb, err := s.AccountRepository.GetByAccountNumber(c, acctNbr)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to get account")
	}
	if accFromDb == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid Account Number/PIN")
	}
	return accFromDb, nil
}

func (s *Service) Transfer(ctx echo.Context, transfer entity.Transfer) (*entity.AccountResponse, *echo.HTTPError) {
	if transfer.FromAccountNumber == "" || transfer.ToAccountNumber == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Account Number is required")
	} else if transfer.FromAccountNumber == transfer.ToAccountNumber {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "From and Destination account number cannot be the same")
	} else if _, err := strconv.Atoi(transfer.FromAccountNumber); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid account")
	}
	fromAccount, err := s.AccountRepository.GetByAccountNumber(ctx, transfer.FromAccountNumber)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to get account")
	}
	toAccount, err := s.AccountRepository.GetByAccountNumber(ctx, transfer.ToAccountNumber)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to get account")
	}
	if fromAccount == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid account")
	} else if toAccount == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid account")
	} else if transfer.Amount <= 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid transfer amount")
	} else if transfer.Amount > 1000 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Maximum amount to transfer is $1000")
	} else if transfer.Amount < 1 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Minimum amount to transfer is $1")
	} else if fromAccount.Balance < transfer.Amount {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Insufficient balance $%0.f", transfer.Amount))
	} else if strings.Trim(transfer.ReferenceNumber, " ") != "" {
		if _, errl := strconv.Atoi(transfer.ReferenceNumber); errl != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid Reference Number")
		}
	}
	fromAccount.Balance -= transfer.Amount
	toAccount.Balance += transfer.Amount
	errl := s.CreateTransactionHistory(ctx, trxEntity.Transaction{
		AccountNumber:           fromAccount.AccountNumber,
		TransferToAccountNumber: toAccount.AccountNumber,
		Amount:                  transfer.Amount,
		Type:                    trxEntity.TYPE_TRANSFER,
	}, nil)
	if err != nil {
		fmt.Printf("Error when creating transaction history: %v \n", errl)
	}
	return fromAccount.ToAccountResponse(), nil
}

func (s *Service) BalanceCheck(ctx echo.Context, acctNbr string) (*entity.AccountResponse, *echo.HTTPError) {
	accFromDb, err := s.GetByAccountNumber(ctx, acctNbr)
	if err != nil {
		return nil, err
	}
	return accFromDb.ToAccountResponse(), nil
}

func (s *Service) Import(c echo.Context, path string) error {
	accounts, err := s.AccountCsvRepository.Import(path)
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		return errors.New("no data imported")
	}
	accMap := make(map[string]*entity.Account)
	for _, ac := range accounts {
		if accMap[ac.AccountNumber] == nil {
			accMap[ac.AccountNumber] = ac
		} else {
			return errors.New("duplicate account number")
		}
	}
	err = s.AccountRepository.BatchInsert(c, accounts)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Insert(ctx echo.Context, account entity.Account) error {
	if strings.TrimSpace(account.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(account.AccountNumber) == "" {
		return errors.New("Account Number is required")
	} else if len(account.AccountNumber) < 6 {
		return errors.New("Account Number should have 6 digits length")
	} else if _, err := strconv.Atoi(account.AccountNumber); err != nil {
		return errors.New("Account Number should only contains numbers")
	}
	accFromDb, _ := s.AccountRepository.GetByAccountNumber(ctx, account.AccountNumber)
	if accFromDb != nil {
		return errors.New("Account number is already exist")
	}
	return s.AccountRepository.Insert(ctx, account)
}

func (s *Service) GetAll(ctx echo.Context) ([]*entity.Account, error) {
	accounts, err := s.AccountRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
