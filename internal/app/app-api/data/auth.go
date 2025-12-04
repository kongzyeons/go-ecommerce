package data

import (
	"app-ecommerce/internal/model"
	"app-ecommerce/pkg/types"
	"strings"
	"time"
)

type AuthRegisterReq struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}

func (obj *AuthRegisterReq) ToReq() AuthRegisterReq {
	obj.Name = strings.TrimSpace(obj.Name)
	obj.Email = strings.TrimSpace(obj.Email)
	obj.Password = strings.TrimSpace(obj.Password)
	return *obj
}

func (obj *AuthRegisterReq) ToUserModel() model.User {
	return model.User{
		Name:         obj.Name,
		Email:        obj.Email,
		Password:     obj.Password,
		RoleID:       1,
		CreatedBy:    "admin",
		CreatedDate:  time.Now().UTC(),
		ModifiedBy:   "admin",
		ModifiedDate: types.NewNullTime(time.Now().UTC()),
	}
}

type AuthLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type AuthLoginRes struct {
	AccToken string `json:"acc_token"`
}

type AuthRefreshTokenRes struct {
	AccToken string `json:"acc_token"`
}
