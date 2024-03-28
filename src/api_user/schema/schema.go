package schema

import (
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
)

type UserRes struct {
	UserID        string `json:"userId"`
	Username      string `json:"username"`
	Mobile        string `json:"mobile"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	DeletedAt     string `json:"deletedAt"`
}

func MapUser2UserRes(user *model.User) *UserRes {
	deletedAt := core.EmptyString
	if user.DeletedAt.Valid {
		deletedAt = user.DeletedAt.Time.String()
	}
	return &UserRes{
		UserID:        user.UserID,
		Username:      user.Username,
		Mobile:        user.Mobile,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt.String(),
		UpdatedAt:     user.UpdatedAt.String(),
		DeletedAt:     deletedAt,
	}
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenRes struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}
