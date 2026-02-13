// handlers/about.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/models"
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
	categories, err := h.DB.GetAllSkillCategories()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load categories")
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return pages.AboutContent(skills, experiences, education, categories).
			Render(c.Request().Context(), c.Response())
	}
	return pages.AboutPage(skills, experiences, education, categories).
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
	project, _ := h.DB.GetRandomProjectForSkill(id) // nil is fine
	return pages.SkillDetail(skill, project).
		Render(c.Request().Context(), c.Response())
}

func (h *AboutHandler) HandleFilteredSkills(c echo.Context) error {
	catIDStr := c.QueryParam("category_id")

	var skills []models.Skill
	var err error

	if catIDStr == "" || catIDStr == "all" {
		skills, err = h.DB.GetAllSkills()
	} else {
		catID, parseErr := strconv.ParseInt(catIDStr, 10, 64)
		if parseErr != nil {
			return c.String(http.StatusBadRequest, "Invalid category ID")
		}
		skills, err = h.DB.GetSkillsByCategoryID(catID)
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load skills")
	}

	return pages.SkillListWithDetail(skills).
		Render(c.Request().Context(), c.Response())
}
