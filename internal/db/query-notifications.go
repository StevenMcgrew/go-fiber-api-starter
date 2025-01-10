package db

import (
	"go-fiber-api-starter/internal/models"

	"github.com/jackc/pgx/v5"
)

func InsertNotification(n *models.Notification) (models.Notification, error) {
	row, err := One(`INSERT INTO notifications (text_content, has_viewed, user_id)
					 VALUES (@textContent, @hasViewed, @userId)
					 RETURNING *;`,
		pgx.NamedArgs{
			"textContent": n.TextContent,
			"hasViewed":   n.HasViewed,
			"userId":      n.UserId},
		&models.Notification{})
	return row, err
}

func GetNotificationById(id uint) (models.Notification, error) {
	row, err := One("SELECT * FROM notifications WHERE id = @id LIMIT 1;",
		pgx.NamedArgs{"id": id},
		&models.Notification{})
	return row, err
}

func GetAllNotificationsForUserId(userId uint) ([]models.Notification, error) {
	rows, err := Many("SELECT * FROM notifications WHERE user_id = @userId;",
		pgx.NamedArgs{"userId": userId},
		&models.Notification{})
	return rows, err
}

func GetNotifications(page uint, perPage uint, query string) ([]models.Notification, string, error) {
	// Query builder
	qb := NewQueryBuilder(
		page,
		perPage,
		query,
		"notifications",
		[]string{
			"id",
			"text_content",
			"has_viewed",
			"user_id",
			"created_at",
		},
	)

	// Build the query string
	queryString, err := qb.Build()
	if err != nil {
		return nil, "", err
	}

	// Run the query
	rows, err := Many(queryString, pgx.NamedArgs{}, &models.Notification{})
	return rows, queryString, err
}

func DeleteNotificationByIds(noteId uint, userId uint) (models.Notification, error) {
	row, err := One(`DELETE FROM notifications
					 WHERE id = @noteId
					 AND user_id = @userId
					 RETURNING *;`,
		pgx.NamedArgs{
			"noteId": noteId,
			"userId": userId},
		&models.Notification{})
	return row, err
}
