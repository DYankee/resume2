package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo, h *SkillHandler) {
	group := app.Group("/skills")

	group.GET("", h.HandleShowSkills)
}
