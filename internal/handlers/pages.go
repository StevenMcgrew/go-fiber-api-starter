package handlers

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func HomePage(c *fiber.Ctx) error {
	data := struct{ ShowLogin bool }{ShowLogin: true}
	filenames := []string{"root-layout", "header"}
	return renderAndSendHTML(c, data, filenames)
}

func EmailVerificationSuccessPage(c *fiber.Ctx) error {
	data := struct{ ShowLogin bool }{ShowLogin: false}
	filenames := []string{"root-layout", "header", "email-verification-success"}
	return renderAndSendHTML(c, data, filenames)
}

func EmailVerificationFailurePage(c *fiber.Ctx) error {
	data := struct{ ShowLogin bool }{ShowLogin: false}
	filenames := []string{"root-layout", "header", "email-verification-failure"}
	return renderAndSendHTML(c, data, filenames)
}

func renderAndSendHTML(c *fiber.Ctx, data any, filenames []string) error {
	// Get views directory
	wd, err := os.Getwd()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when getting working directory", "data": err.Error()})
	}
	viewsDir := wd + "/internal/views"
	// Get template file paths
	paths := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		paths = append(paths, filepath.Join(viewsDir, filename+".html"))
	}
	// Render and send
	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error while parsing HTML templates", "data": err.Error()})
	}
	if tmpl == nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Parsing html templates produced a nil template", "data": ""})
	}
	c.Set("Content-Type", "text/html")
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error while rendering HTML templates", "data": err.Error()})
	}
	return nil
}
