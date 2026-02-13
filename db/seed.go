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
		{"Lua", "Programming Languages", "Lightweight general purpose scripting language. Mainly used as an embedded scripting language in applications. The language that started my programming journey programming turtles for the computer craft mod for minecraft. Now my goto language for writing quick scripts", "", 70},
		{"JavaScript", "Programming Languages", "High level interpreted language that powers the web. Primarily used by the browser to add interactivity to webpages.", "", 40},
		{"Tailwind CSS", "Frontend", "Utility-first CSS framework for rapid UI development. Im a big fan of how tailwind uses utility classes to make it easy to see and change what rules are applied. Plus with components its easy to reuse commonly occurring sets of rules", "", 60},
		{"HTMX", "Frontend", "HTMX is a javascript library for driving hypermedia-driven interactions without a heavy JS framework. It provides a set of custom html attributes which you use to give html elements the ability to trigger an HTTP request. This allows for the creation of a SPA like app with server side state management and html fragments", "", 50},
		{"BubbleTea", "Frontend", "Golang framework based on the elm architecture. Great way to quickly add a UI to command line apps", "", 40},
		{"Echo", "Backend", "High-performance easy to use Go web framework. My goto goLang web framework", "", 60},
		{"SQLite", "Databases", "Lightweight embedded relational database. Great for small applications and prototyping. My goto database for personal projects", "", 70},
		{"MySQL", "Databases", "Open source high performance relational database.", "", 70},
		{"Docker", "Software", "Containerization for consistent application development and deployment. All my websites use docker for easy deployment and management", "", 50},
		{"Git", "Software", "Version control software for managing and maintaining code bases. Git is in my opinion the gold standard for version control. While ", "", 65},
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
		"A simple website I built to host my blog and show off what im working on. Built using the GOTH stack which is Golang, Templ, Tailwind, and HTMX.",
		"", "https://github.com/you/portfolio", "https://example.com",
	)
	p2ID, _ := db.CreateProject(
		"RipR",
		"A TUI program to help with recording and splitting records with audacity.",
		"A basic TUI that uses audacity along with the musicBrains API to find the length of each track, account for the variance between the length of the original and user recording. It then exports each track to the desired folder. UI built using the bubbleTea go framework",
		"", "https://github.com/you/taskcli", "",
	)

	// Link skills to projects via skill_uses
	for _, skillName := range []string{"GoLang", "HTMX", "Tailwind CSS", "Echo", "SQLite", "Docker"} {
		db.AddSkillToProject(skillIDs[skillName], p1ID)
	}
	for _, skillName := range []string{"GoLang", "BubbleTea"} {
		db.AddSkillToProject(skillIDs[skillName], p2ID)
	}

	//Education
	db.CreateEducation("BS Computer Science", "Suny Polytechnic", 3.2, true)
	db.CreateEducation("AAS Computer Information Systems", "Suny Onondaga Community College", 3.83, false)

	// Work Experience
	db.CreateExperience(
		"Fulfillment Center Warehouse Associate", "Amazon",
		"2023-06", "", "Performed audits on employee performance and provided corrective coaching when necessary. Assessed damaged products and determined if and where they should be disposed of.",
	)
	db.CreateExperience(
		"Switcher", "Fed-Ex Ground",
		"2021-09", "2023-06", "Used an electric switcher to move trailers around the yard and out them in their designated spaces",
	)
	db.CreateExperience(
		"Assistant Director", "Chick-Fil-A",
		"2016-11", "2019-09", "Managed front of house staff in a fast paced restaurant environment. Received and handled guest complaints in person and over the phone.",
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
		false,
	)

	log.Println("Seeding complete.")
}
