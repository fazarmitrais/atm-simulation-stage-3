package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/fazarmitrais/atm-simulation-stage-3/lib/responseFormatter"
	"github.com/labstack/echo/v4"
)

func (s *Service) PINValidation(c echo.Context, account entity.Account) *responseFormatter.ResponseFormatter {
	if strings.Trim(account.AccountNumber, " ") == "" {
		return responseFormatter.New(http.StatusBadRequest, "Account Number is required", true)
	} else if strings.Trim(account.PIN, " ") == "" {
		return responseFormatter.New(http.StatusBadRequest, "PIN is required", true)
	} else if len(account.AccountNumber) < 6 {
		return responseFormatter.New(http.StatusBadRequest, "Account Number should have 6 digits length", true)
	} else if len(account.PIN) < 6 {
		return responseFormatter.New(http.StatusBadRequest, "PIN should have 6 digits length", true)
	} else if _, err := strconv.Atoi(account.AccountNumber); err != nil {
		return responseFormatter.New(http.StatusBadRequest, "Account Number should only contains numbers", true)
	} else if _, err := strconv.Atoi(account.PIN); err != nil {
		return responseFormatter.New(http.StatusBadRequest, "PIN should only contains numbers", true)
	}
	accFromDb, err := s.AccountRepository.GetByAccountNumber(c, account.AccountNumber)
	if err != nil {
		return responseFormatter.New(http.StatusInternalServerError, "Failed to get account", true)
	}
	if accFromDb == nil || accFromDb.PIN != account.PIN {
		return responseFormatter.New(http.StatusBadRequest, "Invalid Account Number/PIN", true)
	}
	return nil
}

func (s *Service) Withdraw(ctx echo.Context, accountNumber string, withdrawAmount float64) (*entity.AccountResponse, *responseFormatter.ResponseFormatter) {
	accFromDb, err := s.getAndValidateByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	if withdrawAmount <= 0 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid withdraw amount", true)
	} else if withdrawAmount > 1000 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Maximum amount to withdraw is $1000", true)
	} else if int(withdrawAmount)%10 != 0 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid ammount", true)
	}

	if accFromDb.Balance < withdrawAmount {
		return nil, responseFormatter.New(http.StatusBadRequest, fmt.Sprintf("Insufficient balance $%0.f", withdrawAmount), true)
	}
	accFromDb.Balance -= withdrawAmount
	errl := s.CreateTransactionHistory(ctx, trxEntity.Transaction{
		AccountNumber: accountNumber,
		Amount:        withdrawAmount,
		Type:          trxEntity.TYPE_WITHDRAWAL,
	})
	if err != nil {
		fmt.Printf("Error when creating transaction history: %v \n", errl)
	}
	return accFromDb.ToAccountResponse(), nil
}

func (s *Service) getAndValidateByAccountNumber(c echo.Context, acctNbr string) (*entity.Account, *responseFormatter.ResponseFormatter) {
	if strings.Trim(acctNbr, " ") == "" {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number is required", true)
	} else if len(acctNbr) < 6 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number should have 6 digits length", true)
	} else if _, err := strconv.Atoi(acctNbr); err != nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number should only contains numbers", true)
	}
	accFromDb, err := s.AccountRepository.GetByAccountNumber(c, acctNbr)
	if err != nil {
		return nil, responseFormatter.New(http.StatusInternalServerError, "Failed to get account", true)
	}
	if accFromDb == nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid Account Number/PIN", true)
	}
	return accFromDb, nil
}

func (s *Service) Transfer(ctx echo.Context, transfer entity.Transfer) (*entity.AccountResponse, *responseFormatter.ResponseFormatter) {
	if transfer.FromAccountNumber == "" || transfer.ToAccountNumber == "" {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number is required", true)
	} else if transfer.FromAccountNumber == transfer.ToAccountNumber {
		return nil, responseFormatter.New(http.StatusBadRequest, "From and Destination account number cannot be the same", true)
	} else if _, err := strconv.Atoi(transfer.FromAccountNumber); err != nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid account", true)
	}
	fromAccount, err := s.AccountRepository.GetByAccountNumber(ctx, transfer.FromAccountNumber)
	if err != nil {
		return nil, responseFormatter.New(http.StatusInternalServerError, "Failed to get account", true)
	}
	toAccount, err := s.AccountRepository.GetByAccountNumber(ctx, transfer.ToAccountNumber)
	if err != nil {
		return nil, responseFormatter.New(http.StatusInternalServerError, "Failed to get account", true)
	}
	if fromAccount == nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid account", true)
	} else if toAccount == nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid account", true)
	} else if transfer.Amount <= 0 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid transfer amount", true)
	} else if transfer.Amount > 1000 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Maximum amount to transfer is $1000", true)
	} else if transfer.Amount < 1 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Minimum amount to transfer is $1", true)
	} else if fromAccount.Balance < transfer.Amount {
		return nil, responseFormatter.New(http.StatusBadRequest, fmt.Sprintf("Insufficient balance $%0.f", transfer.Amount), true)
	} else if strings.Trim(transfer.ReferenceNumber, " ") != "" {
		if _, err = strconv.Atoi(transfer.ReferenceNumber); err != nil {
			return nil, responseFormatter.New(http.StatusBadRequest, "Invalid Reference Number", true)
		}
	}
	fromAccount.Balance -= transfer.Amount
	toAccount.Balance += transfer.Amount
	errl := s.CreateTransactionHistory(ctx, trxEntity.Transaction{
		AccountNumber:           fromAccount.AccountNumber,
		TransferToAccountNumber: toAccount.AccountNumber,
		Amount:                  transfer.Amount,
		Type:                    trxEntity.TYPE_TRANSFER,
	})
	if err != nil {
		fmt.Printf("Error when creating transaction history: %v \n", errl)
	}
	return fromAccount.ToAccountResponse(), nil
}

func (s *Service) BalanceCheck(ctx echo.Context, acctNbr string) (*entity.AccountResponse, *responseFormatter.ResponseFormatter) {
	accFromDb, err := s.getAndValidateByAccountNumber(ctx, acctNbr)
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
