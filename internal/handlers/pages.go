package handlers

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/gofiber/fiber/v2"
)

func HomePage(c *fiber.Ctx) error {
	data := struct{ ShowLogin bool }{ShowLogin: true}
	filenames := []string{"root-layout", "header"}
	return renderAndSendHTML(c, data, filenames)
}

func SuccessfullyVerifiedEmailPage(c *fiber.Ctx) error {
	data := struct {
		ShowLogin bool
	}{
		ShowLogin: false,
	}
	filenames := []string{"root-layout", "header", "email-verification-success"}
	return renderAndSendHTML(c, data, filenames)
}

func FailedToVerifyEmailPage(c *fiber.Ctx, failureMessage string) error {
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
		return fiber.NewError(400, "Error parsing query parameter: "+err.Error())
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

func renderAndSendHTML(c *fiber.Ctx, data any, filenames []string) error {
	// Get views directory
	wd, err := os.Getwd()
	if err != nil {
		return fiber.NewError(500, "Error when getting working directory: "+err.Error())
	}
	viewsDir := wd + "/internal/views"

	// Get template file paths
	paths := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		paths = append(paths, filepath.Join(viewsDir, filename+".html"))
	}

	// Create template
	tmpl, err := template.ParseFiles(paths...)
	if err != nil || tmpl == nil {
		return fiber.NewError(500, "Error creating HTML template: "+err.Error())
	}

	// Set Content-Type
	c.Set("Content-Type", "text/html")

	// Render and send
	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		return fiber.NewError(500, "Error rendering and sending HTML template: "+err.Error())
	}
	return nil
}
