package service

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/model"
	"app-ecommerce/internal/repository"
	"app-ecommerce/pkg/db"
	"app-ecommerce/pkg/response"
	"app-ecommerce/pkg/util"
	"app-ecommerce/pkg/validation"
	"context"
	"time"
)

type AuthSvc interface {
	Register(ctx context.Context, req data.AuthRegisterReq) response.Response[any]
	Login(ctx context.Context, req data.AuthLoginReq) response.Response[*data.AuthUserInfo]
}

type authSvc struct {
	repo repository.Repo
	pg   db.PostgresqlDb
}

func NewAuthSvc() AuthSvc {
	return &authSvc{
		repo: repository.NewRpo(),
		pg:   db.NewPostgresqlDb(),
	}
}

func (svc *authSvc) Register(ctx context.Context, req data.AuthRegisterReq) response.Response[any] {
	req = req.CleanReq()
	if valMap := validation.ValidateRequest(req); len(valMap) > 0 {
		return response.ValidationFailed[any](valMap)
	}

	req.Password = util.HashPassword(req.Password)

	count, err := svc.repo.UserRepo.GetCountUnique(model.User{
		Name: req.Name, Email: req.Email,
	})
	if err != nil {
		return response.InternalServerError[any](err, "error get count unique")
	}
	if count > 0 {
		return response.BadRequest[any]("user already exist")
	}

	err = svc.pg.ExecTx(ctx, func(tx db.TX) error {
		_, err = svc.repo.UserRepo.Insert(tx, req.ToUserModel())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response.InternalServerError[any](err, "error exec tx")
	}

	return response.Ok[any](nil)
}

func (svc *authSvc) Login(ctx context.Context, req data.AuthLoginReq) response.Response[*data.AuthUserInfo] {
	req = req.ToReq()
	if valMap := validation.ValidateRequest(req); len(valMap) > 0 {
		return response.ValidationFailed[*data.AuthUserInfo](valMap)
	}

	dataDB, err := svc.repo.UserRepo.GetInfo(model.User{
		Email: req.Email,
	})
	if err != nil {
		return response.InternalServerError[*data.AuthUserInfo](err, "error get unique")
	}
	if dataDB == nil {
		return response.BadRequest[*data.AuthUserInfo]("user not found")
	}

	hashPassword := util.HashPassword(req.Password)

	if dataDB.Password != hashPassword {
		return response.BadRequest[*data.AuthUserInfo]("password not match")
	}

	// accToken, err := jwt.GenToken(jwt.Token{
	// 	UserID:       dataDB.ID,
	// 	UserName:     dataDB.Name,
	// 	Email:        dataDB.Email,
	// 	Role:         dataDB.RoleName.String,
	// 	TimeDulation: time.Hour,
	// })
	// if err != nil {
	// 	return response.InternalServerError[*data.AuthUserInfo](err, "error gen acc token")
	// }

	// refToken, err := jwt.GenToken(jwt.Token{
	// 	UserID:       dataDB.ID,
	// 	UserName:     dataDB.Name,
	// 	Email:        dataDB.Email,
	// 	Role:         dataDB.RoleName.String,
	// 	TimeDulation: 24 * time.Hour,
	// })
	// if err != nil {
	// 	return response.InternalServerError[*data.AuthUserInfo](err, "error gen ref token")
	// }

	res := data.AuthUserInfo{
		UserID:   dataDB.ID,
		UserName: dataDB.Name,
		Email:    dataDB.Email,
		Role:     dataDB.RoleName.String,
		LastPing: time.Now().UTC().Unix(),
	}
	return response.Ok(&res)
}
