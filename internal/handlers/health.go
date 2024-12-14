package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go-fiber-api-starter/internal/db"

	"github.com/gofiber/fiber/v2"
)

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func HealthCheck(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := db.Pool.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("Database is down: %v", err)
		log.Fatalf("Database is down: %v", err) // Log the error and terminate the program
		return c.JSON(stats)
	}

	// Set initial stats
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats
	s := db.Pool.Stat()
	totalConns := s.TotalConns()
	waitCount := s.EmptyAcquireCount()
	maxIdleDestroyCount := s.MaxIdleDestroyCount()
	maxLifetimeDestroyCount := s.MaxLifetimeDestroyCount()
	acquiredConns := s.AcquiredConns()
	idleConns := s.IdleConns()

	// Set stats
	stats["total_connections"] = strconv.Itoa(int(totalConns))
	stats["current_connections"] = strconv.Itoa(int(acquiredConns))
	stats["idle_connections"] = strconv.Itoa(int(idleConns))
	stats["wait_count"] = strconv.Itoa(int(waitCount))
	stats["max_idle_closed"] = strconv.Itoa(int(maxIdleDestroyCount))
	stats["max_lifetime_closed"] = strconv.Itoa(int(maxLifetimeDestroyCount))

	// Evaluate stats to provide a health message
	if totalConns > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if waitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if int32(maxIdleDestroyCount) > (totalConns / 2) {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if int32(maxLifetimeDestroyCount) > (totalConns / 2) {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return c.JSON(stats)
}
