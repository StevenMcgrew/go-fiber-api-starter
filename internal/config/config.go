package config

import (
	"errors"
	"go-fiber-api-starter/internal/handlers"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Fiber server config options: https://docs.gofiber.io/api/fiber#config
var FiberServerConfig = fiber.Config{
	// AppName:                      "",
	// BodyLimit:                    4194304,
	// CaseSensitive:                false,
	// ColorScheme:                  fiber.DefaultColors,
	// CompressedFileSuffix:         fiber.DefaultCompressedFileSuffix,
	// Concurrency:                  262144,
	// DisableDefaultContentType:    false,
	// DisableDefaultDate:           false,
	// DisableHeaderNormalizing:     false,
	// DisableKeepalive:             false,
	// DisablePreParseMultipartForm: false,
	// DisableStartupMessage:        false,
	// ETag:                         false,
	// EnableIPValidation:           false,
	// EnablePrintRoutes:            false,
	// EnableSplittingOnParsers:     false,
	// EnableTrustedProxyCheck:      false,
	ErrorHandler: ServerErrorHandler,
	// GETOnly:           false,
	// IdleTimeout:       time.Duration(),
	// Immutable:         false,
	// JSONDecoder:       json.Unmarshal,
	// JSONEncoder:       json.Marshal,
	// Network:           fiber.NetworkTCP4,
	// PassLocalsToViews: false,
	// Prefork:           false,
	// ProxyHeader:       "",
	// ReadBufferSize:    fiber.DefaultReadBufferSize,
	// ReadTimeout:       time.Duration(),
	// RequestMethods:    fiber.DefaultMethods,
	// ServerHeader:      "",
	// StreamRequestBody: false,
	// StrictRouting:     false,
	// TrustedProxies:    []string,
	// UnescapePath:      false,
	// Views:             nil,
	// ViewsLayout:       "",
	// WriteBufferSize:   fiber.DefaultWriteBufferSize,
	// WriteTimeout:      time.Duration(),
	// XMLEncoder:        xml.Marshal,
}

// Fiber Logger config options: https://docs.gofiber.io/api/middleware/logger
var FiberLoggerConfig = logger.Config{
	Next:          nil,
	Done:          nil,
	Format:        "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	TimeFormat:    "15:04:05",
	TimeZone:      "Local",
	TimeInterval:  500 * time.Millisecond,
	Output:        os.Stdout,
	DisableColors: false,
}

// Fiber CORS config options: https://docs.gofiber.io/api/middleware/cors
var FiberCorsConfig = cors.Config{
	Next:             nil,
	AllowOriginsFunc: nil,
	AllowOrigins:     "*",
	AllowMethods: strings.Join([]string{
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodHead,
		fiber.MethodPut,
		fiber.MethodDelete,
		fiber.MethodPatch,
	}, ","),
	AllowHeaders:     "Accept,Authorization,Content-Type",
	AllowCredentials: false, // credentials require explicit origins
	ExposeHeaders:    "",
	MaxAge:           300,
}

// Fiber Static config options: https://docs.gofiber.io/api/app
var FiberStaticConfig = fiber.Static{
	Compress:       false,
	ByteRange:      false,
	Browse:         false,
	Download:       false,
	Index:          "index.html",
	CacheDuration:  -1,
	MaxAge:         0,
	ModifyResponse: nil,
	Next:           nil,
}

func ServerErrorHandler(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error

	// Handle errors other than fiber.Error
	if !errors.As(err, &fiberErr) {
		return c.Status(500).SendString("Internal Server Error: " + err.Error())
	}

	// Handle fiber.Error
	fiberErr, ok := err.(*fiber.Error)
	if !ok {
		return c.Status(500).SendString("Internal Server Error: " + err.Error())
	}
	if fiberErr.Code == 404 {
		handlerErr := handlers.ErrorPage(c, fiberErr)
		if handlerErr != nil {
			return c.Status(500).SendString("Internal Server Error: " + fiberErr.Error())
		}
	}
	sendJsonErr := c.Status(fiberErr.Code).JSON(fiber.Map{
		"status":  "error",
		"code":    fiberErr.Code,
		"message": fiberErr.Message,
		"error":   fiberErr.Error(),
		"data":    map[string]any{},
		"pagination": map[string]any{
			"page":       0,
			"perPage":    0,
			"totalPages": 0,
			"totalCount": 0,
			"links": map[string]any{
				"self":     "",
				"first":    "",
				"previous": "",
				"next":     "",
				"last":     "",
			},
		},
	})
	if sendJsonErr != nil {
		return c.Status(500).SendString("Internal Server Error: " + fiberErr.Error())
	}

	return nil
}
