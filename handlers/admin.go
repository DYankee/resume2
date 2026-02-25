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

// ── Experience ──────────────────────────────────────
func (h *AdminHandler) HandleAdminExperience(c echo.Context) error {
	experience, err := h.DB.GetAllExperiences()
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to load Experiences",
		)
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return pages.AdminExperienceContent(experience).
			Render(c.Request().Context(), c.Response())
	}
	return pages.AdminExperiencePage(experience).
		Render(c.Request().Context(), c.Response())
}

func (h *AdminHandler) HandleAdminExperienceForm(c echo.Context) error {

	idStr := c.Param("id")
	if idStr == "" {
		// New skill form
		return pages.ExperienceForm(nil).
			Render(c.Request().Context(), c.Response())
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid experience ID")
	}
	experience, err := h.DB.GetExperienceByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Experience not found")
	}
	return pages.ExperienceForm(experience).
		Render(c.Request().Context(), c.Response())
}

func (h *AdminHandler) HandleCreateExperience(c echo.Context) error {
	title := c.FormValue("title")
	company := c.FormValue("company")
	description := c.FormValue("description")
	startDate := c.FormValue("start_date")
	endDate := c.FormValue("end_date")

	if title == "" {
		return c.String(http.StatusBadRequest, "Title required")
	}
	if company == "" {
		return c.String(http.StatusBadRequest, "company required")
	}
	if startDate == "" {
		return c.String(http.StatusBadRequest, "start Date required")
	}

	_, err := h.DB.CreateExperience(
		title, company, startDate, endDate, description,
	)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to create project",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshExperience")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleUpdateExperience(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid experience ID")
	}

	title := c.FormValue("title")
	company := c.FormValue("company")
	description := c.FormValue("description")
	startDate := c.FormValue("start_date")
	endDate := c.FormValue("end_date")

	err = h.DB.UpdateExperience(
		id, title, company, startDate, endDate, description,
	)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to update Experience",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshExperience")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleDeleteExperience(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid experience ID")
	}

	err = h.DB.SoftDeleteExperience(id)
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to delete project",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshExperience")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleAdminExperienceTable(c echo.Context) error {
	experience, err := h.DB.GetAllExperiences()
	if err != nil {
		return c.String(
			http.StatusInternalServerError, "Failed to load experience",
		)
	}
	return pages.ExperienceTable(experience).
		Render(c.Request().Context(), c.Response())
}

// ── Education ─────────────────────────────────────

func (h *AdminHandler) HandleAdminEducation(c echo.Context) error {
	education, err := h.DB.GetAllEducation()
	if err != nil {
		return c.String(
			http.StatusInternalServerError,
			"Failed to load education",
		)
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return pages.AdminEducationContent(education).
			Render(c.Request().Context(), c.Response())
	}
	return pages.AdminEducationPage(education).
		Render(c.Request().Context(), c.Response())
}

func (h *AdminHandler) HandleAdminEducationForm(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return pages.EducationForm(nil).
			Render(c.Request().Context(), c.Response())
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid education ID")
	}
	education, err := h.DB.GetEducationByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Education not found")
	}
	return pages.EducationForm(education).
		Render(c.Request().Context(), c.Response())
}

func (h *AdminHandler) HandleCreateEducation(c echo.Context) error {
	degree := c.FormValue("degree")
	college := c.FormValue("college")
	gpa, _ := strconv.ParseFloat(c.FormValue("gpa"), 64)
	inProgress := c.FormValue("in_progress") == "true"

	if degree == "" {
		return c.String(http.StatusBadRequest, "Degree required")
	}
	if college == "" {
		return c.String(http.StatusBadRequest, "College required")
	}

	err := h.DB.CreateEducation(degree, college, gpa, inProgress)
	if err != nil {
		return c.String(
			http.StatusInternalServerError,
			"Failed to create education",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshEducation")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleUpdateEducation(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid education ID")
	}

	degree := c.FormValue("degree")
	college := c.FormValue("college")
	gpa, _ := strconv.ParseFloat(c.FormValue("gpa"), 64)
	inProgress := c.FormValue("in_progress") == "true"

	err = h.DB.UpdateEducation(id, degree, college, gpa, inProgress)
	if err != nil {
		return c.String(
			http.StatusInternalServerError,
			"Failed to update education",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshEducation")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleDeleteEducation(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid education ID")
	}

	err = h.DB.SoftDeleteEducation(id)
	if err != nil {
		return c.String(
			http.StatusInternalServerError,
			"Failed to delete education",
		)
	}

	c.Response().Header().Set("HX-Trigger", "refreshEducation")
	return c.String(http.StatusOK, "")
}

func (h *AdminHandler) HandleAdminEducationTable(c echo.Context) error {
	education, err := h.DB.GetAllEducation()
	if err != nil {
		return c.String(
			http.StatusInternalServerError,
			"Failed to load education",
		)
	}
	return pages.EducationTable(education).
		Render(c.Request().Context(), c.Response())
}
