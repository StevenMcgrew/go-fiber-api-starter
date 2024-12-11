package config

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

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
