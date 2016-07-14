package middleware

import (
	"github.com/labstack/echo"
	"net/http"
)

type (
	CookieAuthConfig struct {
		Validator CookieAuthValidator
	}

	CookieAuthValidator func(remoteAddress string, sessionId string) bool
)

func CookieAuth(fn CookieAuthValidator) echo.MiddlewareFunc {
	return CookieAuthWithConfig(CookieAuthConfig{fn})
}

func CookieAuthWithConfig(config CookieAuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("sessionId")
			if err == nil {
				sessionId := cookie.Value()
				if config.Validator(c.Request().RemoteAddress(), sessionId) {
					return next(c)
				}
			}

			return c.Redirect(http.StatusFound, "/login")
		}
	}
}
