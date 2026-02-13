package db

import "log"

func (db *DB) Seed() {
	var count int
	db.Conn.QueryRow("SELECT COUNT(*) FROM skills").Scan(&count)
	if count > 0 {
		return
	}
	log.Println("Seeding database...")

	// Skill categories
	categories := []string{"Software", "Programming languages", "Frontend", "Backend"}
	catIDs := make(map[string]int64)
	for _, name := range categories {
		id, _ := db.CreateSkillCategory(name)
		catIDs[name] = id
	}

	// Skills
	skills := []struct {
		name        string
		category    string
		description string
		icon        string
		proficiency int8
	}{
		{"GoLang", "Programming languages", "High level, compiled, general purpose programming language with built in memory management. Has a robust standard library along with an extensive package repository. My goto language for web servers and CLI/TUI applications", "", 60},
		{"C++", "Programming languages", "High level, compiled, general purpose programming language. Has excellent performance but requires manual memory management. My goto language for games and performance critical code", "", 60},
		{"Java", "Programming Languages", "High level, compiles to byte code that runs anywhere with a JVM, garbage collected. My goto language for cross platform development", "", 50},
		{"Python", "Programming Languages", "High level interpreted language, used for scripting, data analysis, and machine learning. My goto language for image processing and machine learning", "", 60},
		{"Lua", "Programming Languages", "Lightweight general purpose scripting language. Mainly used as an embedded scripting language in applications", "", 60},
		{"JavaScript", "Programming Languages", "Frontend and backend development with modern JS/TS.", "", 85},
		{"Tailwind CSS", "Frontend", "Utility-first CSS framework for rapid UI development.", "", 80},
		{"HTMX", "Frontend", "Hypermedia-driven interactions without heavy JS frameworks.", "", 85},
		{"Echo", "Backend", "High-performance Go web framework.", "", 75},
		{"SQLite", "Databases", "Lightweight embedded relational database.", "", 70},
		{"Docker", "Software", "Containerization for consistent development and deployment.", "", 50},
		{"Git", "Software", "Version control software for managing and maintaining code bases.", "", 65},
	}

	skillIDs := make(map[string]int64)
	for _, s := range skills {
		id, _ := db.CreateSkill(s.name, catIDs[s.category], s.description, s.icon, s.proficiency)
		skillIDs[s.name] = id
	}

	// Projects
	p1ID, _ := db.CreateProject(
		"Portfolio Website",
		"My personal portfolio built with the GOTH stack.",
		"A deep dive into building a portfolio with Go, Templ, HTMX, and Tailwind.",
		"", "https://github.com/you/portfolio", "https://example.com",
	)
	p2ID, _ := db.CreateProject(
		"CLI Task Manager",
		"A terminal-based task manager written in Go.",
		"Manage your tasks from the command line with SQLite persistence.",
		"", "https://github.com/you/taskcli", "",
	)

	// Link skills to projects via skill_uses
	for _, skillName := range []string{"Go", "HTMX", "Tailwind CSS", "Echo", "SQLite"} {
		db.AddSkillToProject(skillIDs[skillName], p1ID)
	}
	for _, skillName := range []string{"Go", "SQLite"} {
		db.AddSkillToProject(skillIDs[skillName], p2ID)
	}

	//Education
	db.CreateEducation("BS Computer Science", "Suny Polytechnic", 3.2, true)

	// Experiences
	db.CreateExperience(
		"BS Computer Science", "University of Example",
		"2022-09", "", "Pursuing a BS in Computer Science.",
	)
	db.CreateExperience(
		"Software Intern", "Acme Corp",
		"2024-06", "2024-09", "Built internal tools with Go and React.",
	)

	// Blog posts
	db.CreateBlogPost(
		"Getting Started with the GOTH Stack",
		"getting-started-goth-stack",
		"Learn how to build modern web apps with Go, Templ, HTMX, and Tailwind.",
		`# Getting Started with the GOTH Stack

The GOTH stack combines **Go**, **Templ**, **HTMX**, and **Tailwind CSS** to build fast, server-rendered web applications.

## Why GOTH?

- **Go** gives you a fast, compiled backend.
- **Templ** provides type-safe HTML templates.
- **HTMX** adds dynamic interactions without writing JavaScript.
- **Tailwind** makes styling fast with utility classes.`,
		"Go,HTMX,Tutorial",
		true,
	)

	log.Println("Seeding complete.")
}
