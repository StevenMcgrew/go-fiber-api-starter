package serialization

import (
	"go-fiber-api-starter/internal/models"
)

func ToUserForResponse(user *models.User) *models.UserForResponse {
	return &models.UserForResponse{
		Id:        user.Id,
		Email:     user.Email,
		Username:  user.Username,
		Role:      user.Role,
		Status:    user.Status,
		ImageUrl:  user.ImageUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}
