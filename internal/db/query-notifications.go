package db

import (
	"go-fiber-api-starter/internal/models"

	"github.com/jackc/pgx/v5"
)

func InsertNotification(n *models.Notification) (models.Notification, error) {
	row, err := One(`INSERT INTO notifications (text_content, has_viewed, user_id)
					 VALUES (@textContent, @hasViewed, @userId)
					 RETURNING *;`,
		pgx.NamedArgs{"textContent": n.TextContent, "hasViewed": n.HasViewed, "userId": n.UserId},
		&models.Notification{})
	return row, err
}
