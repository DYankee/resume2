// handlers/projects.go
package handlers

import (
	"net/http"
	"strconv"

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

func (h *ProjectsHandler) HandleProjectExpand(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid project ID")
	}
	project, err := h.DB.GetProjectByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Project not found")
	}
	skills, _ := h.DB.GetSkillsForProject(id)

	pw := pages.ProjectWithSkills{Project: *project, Skills: skills}
	return pages.ProjectCardExpanded(pw).
		Render(c.Request().Context(), c.Response())
}

func (h *ProjectsHandler) HandleProjectCollapse(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid project ID")
	}
	project, err := h.DB.GetProjectByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Project not found")
	}
	skills, _ := h.DB.GetSkillsForProject(id)

	pw := pages.ProjectWithSkills{Project: *project, Skills: skills}
	return pages.ProjectCard(pw).
		Render(c.Request().Context(), c.Response())
}
