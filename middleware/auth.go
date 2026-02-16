// middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/DYankee/resume2/db"
	"github.com/labstack/echo/v4"
)

func RequireAuth(database *db.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil || !database.ValidateSession(cookie.Value) {
				// If HTMX request, tell it to redirect
				if c.Request().Header.Get("HX-Request") == "true" {
					c.Response().Header().Set(
						"HX-Redirect", "/admin/login",
					)
					return c.NoContent(http.StatusUnauthorized)
				}
				return c.Redirect(
					http.StatusSeeOther, "/admin/login",
				)
			}
			return next(c)
		}
	}
}
