package main

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	log.SetLevel(log.LevelTrace)

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(newLogger())
	app.Use(pprof.New())

	app.Get("/", handleIndex)
	app.Get("/about", handleAbout)
	app.Get("/status", newMonitor())
	app.Get("/projects", handleProjects)
	app.Static("/static", "./content/public", fiber.Static{Compress: true})

	app.Listen(":3000")

}

// HANDLERS

func handleIndex(c *fiber.Ctx) error {
	c.Te
	return c.SendString("Hello, World!")
}

func handleAbout(c *fiber.Ctx) error {
	return c.SendString("About")
}

func handleProjects(c *fiber.Ctx) error {
	return c.SendString("Projects")
}

func newMonitor() fiber.Handler {
	return monitor.New(monitor.Config{
		Title:   "resources (jst.dev)",
		Refresh: 3 * time.Second,
		APIOnly: false,
	})
}

// MIDDLEWARE

func newLogger() fiber.Handler {
	return logger.New(logger.Config{
		TimeFormat:    "2006-01-02 15:04:05",
		TimeZone:      "Sweden/Stockholm",
		TimeInterval:  0,
		Output:        nil,
		DisableColors: false,
	})
}
