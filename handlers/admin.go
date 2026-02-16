// handlers/admin.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/templates/pages"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	DB *db.DB
}

// ── Dashboard ─────────────────────────────────────

func (h *AdminHandler) HandleDashboard(c echo.Context) error {
	skills, _ := h.DB.GetAllSkills()
	projects, _ := h.DB.GetAllProjects()
	experiences, _ := h.DB.GetAllExperiences()

	if c.Request().Header.Get("HX-Request") == "true" {
		return pages.AdminDashboardContent(
			len(skills), len(projects), len(experiences),
		).Render(c.Request().Context(), c.Response())
	}
	return pages.AdminDashboardPage(
		len(skills), len(projects), len(experiences),
	).Render(c.Request().Context(), c.Response())
}

// ── Skills ────────────────────────────────────────

func (h *AdminHandler) HandleAdminSkills(c echo.Context) error {
	skills, err := h.DB.GetAllSkills()
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to load skills",
		)
	}
	categories, err := h.DB.GetAllSkillCategories()
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to load categories",
		)
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return pages.AdminSkillsContent(skills, categories).
			Render(c.Request().Context(), c.Response())
	}
	return pages.AdminSkillsPage(skills, categories).
		Render(c.Request().Context(), c.Response())
}

func (h *AdminHandler) HandleAdminSkillForm(c echo.Context) error {
	categories, _ := h.DB.GetAllSkillCategories()

	idStr := c.Param("id")
	if idStr == "" {
		// New skill form
		return pages.SkillForm(nil, categories).
			Render(c.Request().Context(), c.Response())
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid skill ID")
	}
	skill, err := h.DB.GetSkillByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Skill not found")
	}
	return pages.SkillForm(skill, categories).
		Render(c.Request().Context(), c.Response())
}

func (h *AdminHandler) HandleCreateSkill(c echo.Context) error {
	name := c.FormValue("name")
	categoryID, _ := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	description := c.FormValue("description")
	iconURL := c.FormValue("icon_url")
	proficiency, _ := strconv.ParseInt(
		c.FormValue("proficiency"), 10, 8,
	)

	if name == "" || categoryID == 0 {
		return c.String(http.StatusBadRequest, "Name and category required")
	}

	_, err := h.DB.CreateSkill(
		name, categoryID, description, iconURL, int8(proficiency),
	)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to create skill",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshSkills")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleUpdateSkill(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid skill ID")
	}

	name := c.FormValue("name")
	categoryID, _ := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	description := c.FormValue("description")
	iconURL := c.FormValue("icon_url")
	proficiency, _ := strconv.ParseInt(
		c.FormValue("proficiency"), 10, 8,
	)

	err = h.DB.UpdateSkill(
		id, name, categoryID, description, iconURL, int8(proficiency),
	)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to update skill",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshSkills")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleDeleteSkill(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid skill ID")
	}

	err = h.DB.SoftDeleteSkill(id)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to delete skill",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshSkills")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleAdminSkillsTable(c echo.Context) error {
	skills, err := h.DB.GetAllSkills()
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to load skills",
		)
	}
	return pages.SkillsTable(skills).
		Render(c.Request().Context(), c.Response())
}

// ── Skill Categories ──────────────────────────────

func (h *AdminHandler) HandleCreateCategory(c echo.Context) error {
	name := c.FormValue("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "Name required")
	}

	_, err := h.DB.CreateSkillCategory(name)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to create category",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshSkills")
	return c.String(http.StatusOK, "Category created")
}

// ── Projects ──────────────────────────────────────

func (h *AdminHandler) HandleAdminProjects(c echo.Context) error {
	projects, err := h.DB.GetAllProjects()
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to load projects",
		)
	}
	skills, _ := h.DB.GetAllSkills()

	if c.Request().Header.Get("HX-Request") == "true" {
		return pages.AdminProjectsContent(projects, skills).
			Render(c.Request().Context(), c.Response())
	}
	return pages.AdminProjectsPage(projects, skills).
		Render(c.Request().Context(), c.Response())
}

func (h *AdminHandler) HandleAdminProjectForm(c echo.Context) error {
	skills, _ := h.DB.GetAllSkills()

	idStr := c.Param("id")
	if idStr == "" {
		return pages.ProjectForm(nil, skills, nil).
			Render(c.Request().Context(), c.Response())
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid project ID")
	}
	project, err := h.DB.GetProjectByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Project not found")
	}
	projectSkills, _ := h.DB.GetSkillsForProject(id)
	return pages.ProjectForm(project, skills, projectSkills).
		Render(c.Request().Context(), c.Response())
}

func (h *AdminHandler) HandleCreateProject(c echo.Context) error {
	title := c.FormValue("title")
	description := c.FormValue("description")
	longDesc := c.FormValue("long_desc")
	imageURL := c.FormValue("image_url")
	repoURL := c.FormValue("repo_url")
	liveURL := c.FormValue("live_url")

	if title == "" {
		return c.String(http.StatusBadRequest, "Title required")
	}

	projectID, err := h.DB.CreateProject(
		title, description, longDesc, imageURL, repoURL, liveURL,
	)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to create project",
		)
	}

	// Link skills
	form, _ := c.FormParams()
	for _, sid := range form["skill_ids"] {
		skillID, err := strconv.ParseInt(sid, 10, 64)
		if err == nil {
			h.DB.AddSkillToProject(skillID, projectID)
		}
	}

	c.Response().Header().Set("HX-Trigger", "refreshProjects")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleUpdateProject(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid project ID")
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	longDesc := c.FormValue("long_desc")
	imageURL := c.FormValue("image_url")
	repoURL := c.FormValue("repo_url")
	liveURL := c.FormValue("live_url")

	err = h.DB.UpdateProject(
		id, title, description, longDesc, imageURL, repoURL, liveURL,
	)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to update project",
		)
	}

	// Re-link skills: remove old, add new
	oldSkills, _ := h.DB.GetSkillsForProject(id)
	for _, s := range oldSkills {
		h.DB.RemoveSkillFromProject(s.ID, id)
	}
	form, _ := c.FormParams()
	for _, sid := range form["skill_ids"] {
		skillID, err := strconv.ParseInt(sid, 10, 64)
		if err == nil {
			h.DB.AddSkillToProject(skillID, id)
		}
	}

	c.Response().Header().Set("HX-Trigger", "refreshProjects")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleDeleteProject(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid project ID")
	}

	err = h.DB.SoftDeleteProject(id)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to delete project",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshProjects")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleAdminProjectsTable(c echo.Context) error {
	projects, err := h.DB.GetAllProjects()
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to load projects",
		)
	}
	return pages.ProjectsTable(projects).
		Render(c.Request().Context(), c.Response())
}
