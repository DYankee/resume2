package main

import (
	"net/http"

	"github.com/DYankee/resume2/db"
	"github.com/DYankee/resume2/handlers"
	"github.com/DYankee/resume2/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const dbname = "data.db"

func main() {
	app := echo.New()

	app.Static("/", "assets")
	app.Use(middleware.Logger())

	app.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/skill") })

	db, err := db.NewDataStore(dbname)
	if err != nil {
		app.Logger.Fatal(err)
	}

	us := services.NewServicesSkills(services.Skill{}, db)

	h := handlers.New(us)

}
