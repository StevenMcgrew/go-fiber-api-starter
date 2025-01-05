package handlers

import (
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateNotification(c *fiber.Ctx) error {
	// Parse
	notification := &models.Notification{}
	if err := c.BodyParser(notification); err != nil {
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
	}

	// Validate
	length := len(notification.TextContent)
	if length < 1 || length > 500 {
		return fiber.NewError(400, "Text content length is outside of the min/max length")
	}

	// Create (we recreate a notification to only allow setting values that have been validated)
	n := &models.Notification{
		TextContent: notification.TextContent,
		HasViewed:   notification.HasViewed,
		UserId:      notification.UserId,
	}

	// Save
	savedNotification, err := db.InsertNotification(n)
	if err != nil {
		return fiber.NewError(500, "Error saving notification to database: "+err.Error())
	}

	// Respond
	return utils.SendSuccessJSON(c, 201, savedNotification, "Saved new notification")
}

func GetNotification(c *fiber.Ctx) error {
	// Parse
	id, err := c.ParamsInt("noteId")
	if err != nil || id == 0 {
		return fiber.NewError(400, "Error parsing path parameter: "+err.Error())
	}
	noteId := uint(id)

	// Get
	notification, err := db.GetNotificationById(noteId)
	if err != nil {
		return fiber.NewError(500, "Error getting notification from database: "+err.Error())
	}

	// Respond
	return utils.SendSuccessJSON(c, 200, notification, "Retrieved notification from database")
}

func GetAllNotificationsForUser(c *fiber.Ctx) error {
	// Type assert user (the user should be in c.Locals() from AttachUser() middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("user")`)
	}

	// Get
	notifications, err := db.GetAllNotificationsForUserId(user.Id)
	if err != nil {
		return fiber.NewError(500, "Error getting notifications from database: "+err.Error())
	}

	// Respond
	return utils.SendSuccessJSON(c, 200, notifications, "Retrieved all notifications associated with the user id")
}

func GetAllNotifications(c *fiber.Ctx) error {
	return nil
}

func DeleteNotification(c *fiber.Ctx) error {
	// Parse
	id, err := c.ParamsInt("noteId")
	if err != nil || id == 0 {
		return fiber.NewError(400, "Error parsing path parameter: "+err.Error())
	}
	noteId := uint(id)

	// Type assert user (the user should be in c.Locals() from AttachUser() middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("user")`)
	}

	// Delete
	row, err := db.DeleteNotificationByIds(noteId, user.Id)
	if err != nil {
		return fiber.NewError(500, "Error deleting notification from database: "+err.Error())
	}

	// Respond
	data := map[string]any{"id": row.Id}
	return utils.SendSuccessJSON(c, 200, data, ("Deleted the notification"))
}
