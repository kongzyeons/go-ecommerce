package service

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/repository"
	redis_db "app-ecommerce/pkg/redis"
	"app-ecommerce/pkg/response"
	"app-ecommerce/pkg/validation"
	"context"
	"fmt"
	"math"
	"time"
)

type ProductSvc interface {
	GetList(ctx context.Context, req data.ProductGetListReq) response.Response[*data.ProductGetListRes]
}

type productSvc struct {
	repo    repository.Repo
	redisDB redis_db.RedisDB
}

func NewProductSvc() ProductSvc {
	return &productSvc{
		repo:    repository.NewRpo(),
		redisDB: redis_db.NewRedisDB(),
	}
}

func (svc *productSvc) GetList(ctx context.Context, req data.ProductGetListReq) response.Response[*data.ProductGetListRes] {
	req = req.CleanReq()
	if valMap := validation.ValidateRequest(req); len(valMap) > 0 {
		return response.ValidationFailed[*data.ProductGetListRes](valMap)
	}
	if valMap := req.ToVal(); len(valMap) > 0 {
		return response.ValidationFailed[*data.ProductGetListRes](valMap)
	}

	var res data.ProductGetListRes

	initial := req.IsInitial()
	cacheKey := fmt.Sprintf("%s:%s", "productSvc", "getList")
	if initial {
		if svc.redisDB.GetInfo(cacheKey, &res) == nil {
			return response.Ok(&res)
		}
	}

	dataDB, total, err := svc.repo.ProductRepo.GetList(req)
	if err != nil {
		return response.InternalServerError[*data.ProductGetListRes](err, "error get list")
	}

	results := make([]data.ProductGetListResult, len(dataDB))
	for i := range dataDB {
		results[i] = results[i].FromDB(dataDB[i])
	}

	res = data.ProductGetListRes{
		Results:      results,
		TotalPages:   int64(math.Ceil(float64(total) / float64(req.PerPage))),
		TotalResults: total,
		Page:         req.Page,
		PerPage:      req.PerPage,
	}

	if initial {
		svc.redisDB.Set(cacheKey, res, 5*time.Minute)
	}

	return response.Ok(&res)
}
