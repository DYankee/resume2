// handlers/about.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/templates/pages"
	"github.com/labstack/echo/v4"
)

type AboutHandler struct {
	DB *db.DB
}

func (h *AboutHandler) HandleAboutPage(c echo.Context) error {
	skills, err := h.DB.GetAllSkills()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load skills")
	}
	experiences, err := h.DB.GetAllExperiences()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load experiences")
	}
	education, err := h.DB.GetAllEducation()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load education")
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return pages.AboutContent(skills, experiences, education).
			Render(c.Request().Context(), c.Response())
	}
	return pages.AboutPage(skills, experiences, education).
		Render(c.Request().Context(), c.Response())
}

func (h *AboutHandler) HandleSkillDetail(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid skill ID")
	}
	skill, err := h.DB.GetSkillByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Skill not found")
	}
	return pages.SkillDetail(skill).
		Render(c.Request().Context(), c.Response())
}
