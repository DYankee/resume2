package services

import (
	"github.com/DYankee/resume2/db"
)

func NewServicesSkills(s Skill, sStore db.DataStore) *ServicesSkill {
	return &ServicesSkill{
		Skill:      s,
		SkillStore: sStore,
	}
}

type Skill struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	SkillLvl    string `json:"skill_lvl"`
	SkillTypeId string `json:"skill_type_id"`
	Description string `json:"description"`
	IsDeleted   bool   `json:"is_deleted"`
}

type ServicesSkill struct {
	Skill      Skill
	SkillStore db.DataStore
}

func (ss *ServicesSkill) GetAllSkills() ([]Skill, error) {
	query := `SELECT * FROM skills ORDER BY id ASC`

	rows, err := ss.SkillStore.Db.Query(query)
	if err != nil {
		return []Skill{}, err
	}

	defer rows.Close()

	skills := []Skill{}
	for rows.Next() {
		rows.Scan(&ss.Skill)
		skills = append(skills, ss.Skill)
	}
	return skills, nil
}

func (ss *ServicesSkill) GetSkillById(skillId int) (Skill, error) {
	query := `SELECT * FROM skills WHERE id = ?`

	stmt, err := ss.SkillStore.Db.Prepare(query)
	if err != nil {
		return Skill{}, err
	}

	defer stmt.Close()

	ss.Skill.Id = skillId
	err = stmt.QueryRow(ss.Skill.Id).Scan(&ss.Skill)
	if err != nil {
		return Skill{}, err
	}

	return ss.Skill, nil
}
