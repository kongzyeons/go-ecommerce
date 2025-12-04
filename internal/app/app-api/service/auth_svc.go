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
)

type AuthSvc interface {
	Register(ctx context.Context, req data.AuthRegisterReq) response.Response[any]
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
	req = req.ToReq()
	if valMap := validation.ValidateRequest(req); len(valMap) > 0 {
		return response.ValidationFailed[any](valMap)
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return response.InternalServerError[any](err, "error hash password")
	}
	req.Password = hashPassword

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
