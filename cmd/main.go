package main

import (
	"net/http"

	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/handlers"
	"github.com/DYankee/resume2/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const dbname = "../data.db"

func main() {
	app := echo.New()

	app.Static("/assets", "../assets")
	app.Use(middleware.Logger())

	app.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/skills") })

	db, err := db.NewDataStore(dbname)
	if err != nil {
		app.Logger.Fatal(err)
	}

	ss := services.NewServicesSkills(services.Skill{}, db)

	h := handlers.New(ss)

	handlers.SetupRoutes(app, h)

	app.Logger.Fatal(app.Start(":8080"))
}
