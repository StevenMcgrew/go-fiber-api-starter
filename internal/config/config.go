package config

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

var (
	LoginDuration       = 7 * 24 * time.Hour // one week
	VerifyEmailDuration = 15 * time.Minute   // fifteen minutes
	API_BASE_URL        = os.Getenv("API_BASE_URL")
	API_SECRET          = os.Getenv("API_SECRET")
	API_PORT            = os.Getenv("API_PORT")
	API_ENV             = os.Getenv("API_ENV")
	DB_HOST             = os.Getenv("DB_HOST")
	DB_PORT             = os.Getenv("DB_PORT")
	DB_DATABASE         = os.Getenv("DB_DATABASE")
	DB_USERNAME         = os.Getenv("DB_USERNAME")
	DB_PASSWORD         = os.Getenv("DB_PASSWORD")
	DB_SCHEMA           = os.Getenv("DB_SCHEMA")
	DB_URL              = os.Getenv("DB_URL")
	EMAIL_FROM          = os.Getenv("EMAIL_FROM")
	EMAIL_HOST          = os.Getenv("EMAIL_HOST")
	EMAIL_PORT          = os.Getenv("EMAIL_PORT")
	EMAIL_USERNAME      = os.Getenv("EMAIL_USERNAME")
	EMAIL_APP_PASSWORD  = os.Getenv("EMAIL_APP_PASSWORD")
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
	// Views: NewViewEngine(),
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
	MaxAge:           7200,
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

// func NewViewEngine() *html.Engine {
// 	// Get views directory
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		log.Fatal("Error getting working directory during NewViewEngine() function:", err.Error())
// 	}
// 	viewsDir := wd + "/internal/views"

// 	// Create view engine
// 	engine := html.New(viewsDir, ".html")
// 	return engine
// }

func ServerErrorHandler(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error

	// Type assert fiber.Error
	fiberErr, ok := err.(*fiber.Error)
	if !ok {
		// Handle errors other than fiber.Error
		return c.Status(500).SendString("Internal Server Error: " + err.Error())
	}

	// Handle fiber.Error's
	sendJsonErr := c.Status(fiberErr.Code).JSON(fiber.Map{
		"status":  "error",
		"code":    fiberErr.Code,
		"message": fiberErr.Message,
		"error":   fiberErr.Error(),
		"data":    nil,
	})
	if sendJsonErr != nil {
		return c.Status(500).SendString("Internal Server Error: " + fiberErr.Error())
	}
	return nil
}
