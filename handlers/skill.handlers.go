package handlers

import (
	"github.com/DYankee/resume2/services"
	"github.com/DYankee/resume2/views/skill"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type SkillService interface {
	GetAllSkills() ([]services.Skill, error)
	GetSkillById(id int) (services.Skill, error)
}

type SkillHandler struct {
	SkillService SkillService
}

func New(ss SkillService) *SkillHandler {
	return &SkillHandler{
		SkillService: ss,
	}
}

func (sh *SkillHandler) HandleShowSkills(c echo.Context) error {
	skillData, err := sh.SkillService.GetAllSkills()
	if err != nil {
		return err
	}

	si := skill.ShowIndex("test", skill.Show(skillData))

	return sh.View(c, si)
}

func (sh *SkillHandler) View(c echo.Context, tmp templ.Component) error {
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	return tmp.Render(c.Request().Context(), c.Response().Writer)
}
