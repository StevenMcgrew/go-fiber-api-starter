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
	data := struct {
		ShowLogin bool
	}{
		ShowLogin: false,
	}
	filenames := []string{"root-layout", "header", "email-verification-success"}
	return renderAndSendHTML(c, data, filenames)
}

func EmailVerificationFailurePage(c *fiber.Ctx, failureMessage string) error {
	data := struct {
		ShowLogin      bool
		FailureMessage string
	}{
		ShowLogin:      false,
		FailureMessage: failureMessage,
	}
	filenames := []string{"root-layout", "header", "email-verification-failure"}
	c.Status(400)
	return renderAndSendHTML(c, data, filenames)
}

func ResetPasswordPage(c *fiber.Ctx) error {
	// Query params
	type queryParams struct {
		Token string `query:"token"`
	}
	qParams := &queryParams{}

	// Parse query params
	if err := c.QueryParser(qParams); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error parsing 'token' query param",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	data := struct {
		ShowLogin bool
		Token     string
	}{
		ShowLogin: false,
		Token:     qParams.Token,
	}
	filenames := []string{"root-layout", "header", "password-reset"}
	return renderAndSendHTML(c, data, filenames)
}

func ErrorPage(c *fiber.Ctx, err *fiber.Error) error {
	data := struct {
		ShowLogin    bool
		StatusCode   int
		ErrorMessage string
		ErrorDetails string
	}{
		ShowLogin:    false,
		StatusCode:   err.Code,
		ErrorMessage: err.Message,
		ErrorDetails: err.Error(),
	}
	filenames := []string{"root-layout", "header", "error"}
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
