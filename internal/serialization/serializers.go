package serialization

import (
	"go-fiber-api-starter/internal/models"
)

func UserResponse(u *models.User) *models.UserResponse {
	return &models.UserResponse{
		Id:        u.Id,
		Email:     u.Email,
		Username:  u.Username,
		Role:      u.Role,
		Status:    u.Status,
		ImageUrl:  u.ImageUrl,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}
