package handlers

import (
	"fmt"
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/utils"
	"math"

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

// ?page=5&per_page=20&query=where.has_viewed.eq.false.orderby.id.desc
func GetAllNotifications(c *fiber.Ctx) error {
	// Expected query parameters
	type queryParams struct {
		Page    uint   `query:"page"`
		PerPage uint   `query:"per_page"`
		Query   string `query:"query"`
	}
	qParams := &queryParams{}

	// Parse (this also unescapes any escape sequences from the query param values)
	if err := c.QueryParser(qParams); err != nil {
		return fiber.NewError(400, "Error parsing query parameters: "+err.Error())
	}

	// Set simpler var names
	page := qParams.Page
	perPage := qParams.PerPage
	query := qParams.Query

	// Get row rowCount
	rowCount, err := db.GetRowCount("notificaTions")
	if err != nil {
		return fiber.NewError(500, "Error getting row count: "+err.Error())
	}
	if rowCount == 0 {
		return fiber.NewError(400, "No records were found in the database")
	}

	// Validate
	if page < 1 {
		return fiber.NewError(400, "Page number must be 1 or greater")
	}
	floatPageCount := math.Ceil(float64(rowCount) / float64(perPage))
	pageCount := uint(floatPageCount)
	if page > pageCount {
		return fiber.NewError(400, "The page number requested is larger than the total number of pages")
	}

	// Get notifications
	notifications, sql, err := db.GetNotifications(page, perPage, query)
	if err != nil {
		return fiber.NewError(500, "Error getting notifications from database: "+err.Error())
	}

	// Create pagination data for response
	pre := "/api/v1/notifications"
	selfLink := fmt.Sprintf("%s?page=%d&per_page=%d&query=%s", pre, page, perPage, query)
	firstLink := fmt.Sprintf("%s?page=%d&per_page=%d&query=%s", pre, 1, perPage, query)
	previousLink := fmt.Sprintf("%s?page=%d&per_page=%d&query=%s", pre, page-1, perPage, query)
	if page == 1 {
		previousLink = ""
	}
	nextLink := fmt.Sprintf("%s?page=%d&per_page=%d&query=%s", pre, page+1, perPage, query)
	if page == pageCount {
		nextLink = ""
	}
	lastLink := fmt.Sprintf("%s?page=%d&per_page=%d&query=%s", pre, pageCount, perPage, query)
	pageData := &models.Pagination{
		Page:         page,
		PerPage:      perPage,
		TotalPages:   pageCount,
		TotalCount:   rowCount,
		SelfLink:     selfLink,
		FirstLink:    firstLink,
		PreviousLink: previousLink,
		NextLink:     nextLink,
		LastLink:     lastLink,
	}

	// Respond
	return utils.SendPaginationJSON(c, notifications, pageData, sql)
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
