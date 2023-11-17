package cookies

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetCookie(c echo.Context, accountNumber string) {
	config := new(http.Cookie)
	config.Name = "AccountNumber"
	config.Value = accountNumber
	config.Expires = time.Now().Add(1 * time.Hour)
	config.Path = "/"
	c.SetCookie(config)
}

func DeleteCookie(c echo.Context) *echo.HTTPError {
	cookie, _ := c.Cookie("AccountNumber")
	if cookie != nil {
		cookie.Expires = time.Now()
	}
	return nil
}

func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie("AccountNumber")
		if err != nil {
			response := make(map[string]string)
			response["message"] = "Please login first"
			return c.Render(http.StatusUnauthorized, "login.html", response)
		}
		return next(c)
	}
}

func GetAccountNumberFromCookie(c echo.Context) string {
	co, err := c.Cookie("AccountNumber")
	if err != nil || co == nil {
		response := make(map[string]string)
		response["message"] = "Please login first"
		c.Render(http.StatusUnauthorized, "login.html", response)
	}
	return co.Value
}
