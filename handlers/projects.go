// handlers/projects.go
package handlers

import (
	"net/http"

	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/templates/pages"
	"github.com/labstack/echo/v4"
)

type ProjectsHandler struct {
	DB *db.DB
}

func (h *ProjectsHandler) HandleProjectsPage(c echo.Context) error {
	projects, err := h.DB.GetAllProjects()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load projects")
	}

	var items []pages.ProjectWithSkills
	for _, p := range projects {
		skills, err := h.DB.GetSkillsForProject(p.ID)
		if err != nil {
			skills = nil // degrade gracefully
		}
		items = append(items, pages.ProjectWithSkills{
			Project: p,
			Skills:  skills,
		})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return pages.ProjectsContent(items).
			Render(c.Request().Context(), c.Response())
	}
	return pages.ProjectsPage(items).
		Render(c.Request().Context(), c.Response())
}
