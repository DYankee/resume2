// main.go
package main

import (
	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/handlers"
	customMw "github.com/DYankee/resume2/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database := db.New("data/portfolio.db")
	database.Seed()
	defer database.Conn.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(
		middleware.GzipConfig{Level: 5},
	))

	e.Static("/static", "static")

	aboutH := &handlers.AboutHandler{DB: database}
	projectsH := &handlers.ProjectsHandler{DB: database}
	adminH := &handlers.AdminHandler{DB: database}
	authH := &handlers.AuthHandler{DB: database}

	// Public pages
	e.GET("/", aboutH.HandleAboutPage)
	e.GET("/projects", projectsH.HandleProjectsPage)

	// Public HTMX endpoints
	e.GET("/api/skills", aboutH.HandleFilteredSkills)
	e.GET("/api/skills/:id", aboutH.HandleSkillDetail)
	e.GET(
		"/api/projects/:id/expand", projectsH.HandleProjectExpand,
	)
	e.GET(
		"/api/projects/:id/collapse", projectsH.HandleProjectCollapse,
	)

	// Auth routes (public, no middleware)
	e.GET("/admin/login", authH.HandleLoginPage)
	e.POST("/admin/login", authH.HandleLogin)
	e.POST("/admin/logout", authH.HandleLogout)

	// Protected Admin pages
	admin := e.Group("/admin")
	admin.Use(customMw.RequireAuth(database))

	admin.GET("", adminH.HandleDashboard)
	admin.GET("/skills", adminH.HandleAdminSkills)
	admin.GET("/skills/table", adminH.HandleAdminSkillsTable)
	admin.GET("/skills/new", adminH.HandleAdminSkillForm)
	admin.GET("/skills/:id/edit", adminH.HandleAdminSkillForm)
	admin.POST("/skills", adminH.HandleCreateSkill)
	admin.PUT("/skills/:id", adminH.HandleUpdateSkill)
	admin.DELETE("/skills/:id", adminH.HandleDeleteSkill)

	admin.POST("/categories", adminH.HandleCreateCategory)

	admin.GET("/projects", adminH.HandleAdminProjects)
	admin.GET("/projects/table", adminH.HandleAdminProjectsTable)
	admin.GET("/projects/new", adminH.HandleAdminProjectForm)
	admin.GET("/projects/:id/edit", adminH.HandleAdminProjectForm)
	admin.POST("/projects", adminH.HandleCreateProject)
	admin.PUT("/projects/:id", adminH.HandleUpdateProject)
	admin.DELETE("/projects/:id", adminH.HandleDeleteProject)

	e.Logger.Fatal(e.Start(":8080"))
}
