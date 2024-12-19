package config

import (
	"os"
	"strings"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
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

// Fiber JWT config options: https://docs.gofiber.io/contrib/jwt_v1.x.x/jwt/#config
var FiberJwtConfig = jwtware.Config{
	Filter:         nil,
	SuccessHandler: nil,
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if err.Error() == "Missing or malformed JWT" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
		}
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	},
	SigningKey:  jwtware.SigningKey{Key: []byte(os.Getenv("SECRET"))},
	SigningKeys: nil,
	ContextKey:  "payload",
	Claims:      jwt.MapClaims{},
	TokenLookup: "header:Authorization",
	AuthScheme:  "Bearer",
	KeyFunc:     nil,
	JWKSetURLs:  nil,
}
