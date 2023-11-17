package menu

import (
	"fmt"
	"time"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/fazarmitrais/atm-simulation-stage-3/service"
	"github.com/labstack/echo/v4"
)

type Menu struct {
	service       *service.Service
	accountNumber *string
}

func NewMenu(s *service.Service) *Menu {
	return &Menu{service: s}
}

func (m *Menu) Start(c echo.Context) {
	for {
		m.LoginScreen(c)
		m.TransactionScreen(c)
	}
}

func (m *Menu) LoginScreen(c echo.Context) {
	for m.accountNumber == nil {
		var accountNumber, pin string
		fmt.Println()
		fmt.Println("___Welcome To ATM Simulation___")
		fmt.Print("Enter Account Number: ")
		fmt.Scanln(&accountNumber)
		fmt.Print("Enter PIN: ")
		fmt.Scanln(&pin)
		resp := m.service.PINValidation(c, entity.Account{
			AccountNumber: accountNumber,
			PIN:           pin,
		})
		if resp != nil {
			fmt.Println(resp.Message)
		} else {
			m.accountNumber = &accountNumber
		}
	}
}

func (m *Menu) TransactionScreen(c echo.Context) {
	for m.accountNumber != nil {
		fmt.Println()
		fmt.Println("___Transactions___")
		fmt.Println("1. Withdraw")
		fmt.Println("2. Fund Transfer")
		fmt.Println("3. Balance Check")
		fmt.Println("4. Get Transaction History")
		fmt.Println("5. Exit")
		var input int
		fmt.Scanln(&input)
		switch input {
		case 1:
			m.WithdrawScreen(c)
		case 2:
			m.TransferScreen(c)
		case 3:
			m.BalanceCheck(c)
		case 4:
			m.GetLastTransaction(c)
		case 5:
			m.Exit()
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func (m *Menu) GetLastTransaction(c echo.Context) {
	fmt.Println()
	fmt.Println("___Last Transaction___")
	trxs, err := m.service.GetLastTransaction(c, *m.accountNumber, nil, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, t := range trxs {
		trfTo := ""
		if t.Type == trxEntity.TYPE_TRANSFER {
			trfTo = fmt.Sprintf("Transfer to : %s, ", t.TransferToAccountNumber)
		}
		fmt.Printf("Transaction type : %s, %sAmount: $%0.f, Transaction date : %s \n", t.Type, trfTo, t.Amount, t.Date.Format("2006-01-02 15:04 AM"))
	}
	fmt.Println()
	fmt.Println("Press any key to return to the menu screen")
	fmt.Scanln()
}

func (m *Menu) BalanceCheck(c echo.Context) {
	acc, resp := m.service.BalanceCheck(c, *m.accountNumber)
	if resp != nil {
		fmt.Println(resp.Message)
	} else {
		fmt.Printf("Your balance is $%0.f \n", acc.Balance)
		fmt.Println("Press enter to return to main menu")
		fmt.Scanln()
	}
}

func (m *Menu) WithdrawScreen(c echo.Context) {
	for m.accountNumber != nil {
		fmt.Println()
		fmt.Println("___Withdrawal___")
		fmt.Println("1. $10")
		fmt.Println("2. $50")
		fmt.Println("3. $100")
		fmt.Println("4. Other")
		fmt.Println("5. Back")
		var input int
		fmt.Scanln(&input)
		var amount float64
		switch input {
		case 5:
			return
		case 1:
			amount = 10
		case 2:
			amount = 50
		case 3:
			amount = 100
		case 4:
			amount = m.OtherWithdrawScreen(c)
		default:
			fmt.Println("Invalid Input")
		}
		acc, resp := m.service.Withdraw(c, *m.accountNumber, amount)
		if resp != nil {
			fmt.Println(resp.Message)
		} else {
			fmt.Println()
			fmt.Println("___Summary___")
			fmt.Printf("Date : %s\n", time.Now().Format("2006-01-02 15:04 AM")) //TODO
			fmt.Printf("Withdraw : $%0.f\n", amount)
			fmt.Printf("Balance : $%0.f\n", acc.Balance)
			fmt.Println()
			fmt.Println("1. Transaction")
			fmt.Println("2. Exit")
			var input int
			fmt.Scanln(&input)
			switch input {
			case 1:
				return
			case 2:
				m.Exit()
			default:
				fmt.Println("Invalid Input")
			}
		}
	}
}

func (m *Menu) Exit() {
	m.accountNumber = nil
}

func (m *Menu) OtherWithdrawScreen(c echo.Context) float64 {
	fmt.Println()
	fmt.Println("___Other Withdraw___")
	fmt.Print("Enter amount to withdraw: ")
	var amount float64
	fmt.Scanln(&amount)
	return amount
}

func (m *Menu) TransferScreen(c echo.Context) {
	fmt.Println("Please enter destination account and")
	fmt.Println("press enter to continue or")
	fmt.Print("press cancel (Esc) to go back to Transaction: ")
	var input string
	fmt.Scanln(&input)
	if input == "cancel" {
		return
	}
	destAcctNbr := input
	fmt.Println("Please enter transfer amount and press enter to continue or ")
	fmt.Print("press enter to go back to Transaction: ")
	var amount float64
	fmt.Scanln(&amount)
	refNumber := "123456" // TODO
	fmt.Printf("Reference Number: %s \n", refNumber)
	fmt.Print("press enter to continue or press enter to go back to Transaction: ")
	fmt.Scanln()
	fmt.Println("___Transfer Confirmation___")
	fmt.Printf("Destination Account : %s \n", destAcctNbr)
	fmt.Printf("Transfer Amount : $%0.f \n", amount)
	fmt.Printf("Reference Number : %s \n", refNumber)
	fmt.Println()
	fmt.Println("1. Confirm Trx")
	fmt.Println("2. Cancel Trx")
	var input2 int
	fmt.Scanln(&input2)
	if input2 == 2 {
		return
	} else if input2 == 1 {
		acc, resp := m.service.Transfer(c, entity.Transfer{
			FromAccountNumber: *m.accountNumber,
			ToAccountNumber:   destAcctNbr,
			ReferenceNumber:   refNumber,
			Amount:            amount,
		})
		if resp != nil {
			fmt.Println(resp.Message)
		} else {
			fmt.Println()
			fmt.Println("___Fund Transfer Summary___")
			fmt.Printf("Destination Account : %s \n", destAcctNbr)
			fmt.Printf("Transfer Amount : $%0.f \n", amount)
			fmt.Printf("Reference Number : %s \n", refNumber)
			fmt.Printf("Balance : $%0.f\n", acc.Balance)
			fmt.Println()
			fmt.Println("1. Transaction")
			fmt.Println("2. Exit")
			var input int
			fmt.Scanln(&input)
			switch input {
			case 1:
				return
			case 2:
				m.Exit()
			default:
				fmt.Println("Invalid Input")
			}
		}
	}
}
