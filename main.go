// main.go
package main

import (
	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database := db.New("db/portfolio.db")
	database.Seed()
	defer database.Conn.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	e.Static("/static", "static")

	aboutH := &handlers.AboutHandler{DB: database}
	projectsH := &handlers.ProjectsHandler{DB: database}
	//blogH := &handlers.BlogHandler{DB: database}

	// Pages
	e.GET("/", aboutH.HandleAboutPage)
	e.GET("/projects", projectsH.HandleProjectsPage)
	//e.GET("/blog", blogH.HandleBlogPage)
	//e.GET("/blog/:slug", blogH.HandleBlogPost)

	// HTMX API endpoints
	e.GET("/api/skills", aboutH.HandleFilteredSkills)
	e.GET("/api/skills/:id", aboutH.HandleSkillDetail)
	e.GET("/api/projects/:id/expand", projectsH.HandleProjectExpand)
	e.GET("/api/projects/:id/collapse", projectsH.HandleProjectCollapse)

	e.Logger.Fatal(e.Start(":8080"))
}
