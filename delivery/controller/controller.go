package Controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fazarmitrais/atm-simulation-stage-3/cookies"
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/fazarmitrais/atm-simulation-stage-3/service"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	service *service.Service
}

type ResponseFormatter struct {
	IsError bool   `json:"isError"`
	Mesage  string `json:"message"`
}

func New(svc *service.Service) *Controller {
	return &Controller{service: svc}
}

func (re *Controller) Register(e *echo.Echo) {
	auth := e.Group("/api/v1/auth")
	auth.GET("/login", re.PINValidation)
	auth.POST("/login", re.PINValidation)
	auth.GET("/logout", re.Logout)
	account := e.Group("/api/v1/account", cookies.Authorize)
	account.GET("", re.Accounts)
	account.POST("/import", re.Import)
	account.GET("/create", re.Create)
	account.POST("/create", re.Create)
	withdraw := e.Group("/api/v1/withdraw", cookies.Authorize)
	withdraw.GET("", re.Withdraw)
	withdraw.POST("", re.Withdraw)
	withdraw.GET("/summary", re.WithdrawSummary)
}

var response = make(map[string]interface{})

func (re *Controller) Withdraw(c echo.Context) error {
	var statusCode = http.StatusOK
	if c.Request().Method == http.MethodPost {
		accNbr := cookies.GetAccountNumberFromCookie(c)
		type transferAmount struct {
			Amount float64 `json:"amount"`
		}
		amt := transferAmount{}
		c.Bind(&amt)
		_, resp := re.service.Withdraw(c, accNbr, amt.Amount)
		if resp != nil {
			statusCode = resp.Code
			response["message"] = resp.Message
		}
	}
	return c.Render(statusCode, "withdraw.html", response)
}

func (re *Controller) WithdrawSummary(c echo.Context) error {
	accNbr := cookies.GetAccountNumberFromCookie(c)
	tp := string(trxEntity.TYPE_WITHDRAWAL)
	var statusCode int = http.StatusOK
	trx, err := re.service.GetLastTransaction(c, accNbr, &tp, 1)
	if err != nil {
		statusCode = err.Code
		response["message"] = err.Error()
	} else if len(trx) > 0 {
		acc, err := re.service.GetByAccountNumber(c, accNbr)
		if err != nil {
			statusCode = err.Code
			response["message"] = err.Error()
		}
		response["Date"] = trx[0].Date.Format("2006-01-02 15:04")
		response["Withdraw"] = fmt.Sprintf("%0.f", trx[0].Amount)
		response["Balance"] = fmt.Sprintf("%0.f", acc.Balance)
	}
	return c.Render(statusCode, "withdrawSummary.html", response)
}

func (re *Controller) Logout(c echo.Context) error {
	err := cookies.DeleteCookie(c)
	if err != nil {
		response["message"] = err.Message
		return c.Render(err.Code, "index.html", response)
	}
	response["message"] = "Successfully logged out"
	return c.Render(http.StatusOK, "login.html", response)
}

func (re *Controller) Create(c echo.Context) error {
	var err error
	if c.Request().Method == http.MethodPost {
		var account entity.Account
		err = c.Bind(&account)
		if err == nil {
			err = re.service.Insert(c, account)
		}
		if err != nil {
			response["message"] = err.Error()
		} else {
			response["message"] = "Success"
		}
	}
	return c.Render(http.StatusOK, "accountForm.html", response)
}

func (re *Controller) Accounts(c echo.Context) error {
	acc, err := re.service.GetAll(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.Render(http.StatusOK, "index.html", acc)
}

func (re *Controller) Import(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error parsing form")
	}
	files := form.File["fileInput"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Error retrieving file information")
	}
	file := files[0]
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error opening file")
	}
	defer src.Close()

	dst, err := os.Create("uploads/" + file.Filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creating file")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "Error copying file")
	}

	err = re.service.Import(c, "uploads/"+file.Filename)
	if err != nil {
		return err
	}
	return re.Accounts(c)
}

func (re *Controller) PINValidation(c echo.Context) error {
	var statusCode int = http.StatusOK
	if c.Request().Method == http.MethodPost {
		var acc entity.Account
		c.Bind(&acc)
		err := re.service.PINValidation(c, acc)
		if err != nil {
			statusCode = err.Code
			response["message"] = err.Message
		}
		cookies.SetCookie(c, acc.AccountNumber)
		co, _ := c.Cookie("AccountNumber")
		log.Println(co)
	}
	return c.Render(statusCode, "login.html", response)
}

/*
func (re *Controller) Register(e *echo.Echo) {

	g := e.Group("/api/v1/account")
	g.POST("/validate", re.PINValidation)
	g.POST("/withdraw", re.Withdraw)
	g.POST("/transfer", re.Transfer)
	g.GET("/balance", re.BalanceCheck)
	g.GET("/exit", re.Exit)
}

func (re *Controller) BalanceCheck(c echo.Context) error {
	acct, resp := re.service.BalanceCheck(c, cok.Values["acctNbr"].(string))
	if resp != nil {
		resp.ReturnAsJson(w)
		return
	}
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(acct)
	return
}

func (re *Controller) Exit(c echo.Context) error {
	cookieStore, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	cookieStore.Values["authenticated"] = false
	cookieStore.Values["acctNbr"] = nil
	cookieStore.Save(r, w)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(echo.NewHTTPError(http.StatusOK, "Logout success", false))
}

func (re *Controller) Transfer(c echo.Context) error {
	cookieStore, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	var transfer entity.Transfer
	err = json.Unmarshal(b, &transfer)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	transfer.FromAccountNumber = cookieStore.Values["acctNbr"].(string)
	acc, resp := re.service.Transfer(r.Context(), transfer)
	if resp != nil {
		resp.ReturnAsJson(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(acc)
}
*/
