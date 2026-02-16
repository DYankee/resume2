// handlers/auth.go
package handlers

import (
	"crypto/subtle"
	"net/http"
	"os"
	"time"

	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/templates/pages"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	DB *db.DB
}

func (h *AuthHandler) HandleLoginPage(c echo.Context) error {
	return pages.LoginPage("").
		Render(c.Request().Context(), c.Response())
}

func (h *AuthHandler) HandleLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	expectedUser := os.Getenv("ADMIN_USER")
	expectedPass := os.Getenv("ADMIN_PASS")

	userMatch := subtle.ConstantTimeCompare(
		[]byte(username), []byte(expectedUser),
	) == 1
	passMatch := subtle.ConstantTimeCompare(
		[]byte(password), []byte(expectedPass),
	) == 1

	if !userMatch || !passMatch {
		// Return the form with an error (works with HTMX)
		return pages.LoginForm("Invalid username or password").
			Render(c.Request().Context(), c.Response())
	}

	// Create a session lasting 7 days
	token, err := h.DB.CreateSession(7 * 24 * time.Hour)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Session error",
		)
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/admin",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	}
	c.SetCookie(cookie)

	// Tell HTMX to redirect
	c.Response().Header().Set("HX-Redirect", "/admin")
	return c.NoContent(http.StatusOK)
}

func (h *AuthHandler) HandleLogout(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err == nil {
		h.DB.DeleteSession(cookie.Value)
	}

	// Clear the cookie
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/admin",
		HttpOnly: true,
		MaxAge:   -1,
	})

	return c.Redirect(http.StatusSeeOther, "/admin/login")
}
