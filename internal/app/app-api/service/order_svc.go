package service

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/events"
	"app-ecommerce/internal/meta"
	"app-ecommerce/internal/model"
	"app-ecommerce/internal/repository"
	"app-ecommerce/pkg/db"
	"app-ecommerce/pkg/kafka"
	redis_db "app-ecommerce/pkg/redis"
	"app-ecommerce/pkg/response"
	"app-ecommerce/pkg/validation"
	"context"
	"fmt"
	"math"
	"strings"
	"time"
)

type OrderSvc interface {
	Create(ctx context.Context, req data.OrderCreateReq) response.Response[*data.OrderCreateRes]
	GetHistory(ctx context.Context, req data.OrderGetHistoryReq) response.Response[*data.OrderGetHistoryRes]
	Delete(ctx context.Context, req data.OrderDeleteReq) response.Response[*data.OrderDeleteRes]
	Confirm(ctx context.Context, req data.OrderConfirmReq) response.Response[*data.OrderConfirmRes]
	Shipping(ctx context.Context, req data.OrderShippingReq) response.Response[*data.OrderShippingRes]
	Completed(ctx context.Context, req data.OrderCompletedReq) response.Response[*data.OrderCompletedRes]
	Cancel(ctx context.Context, req data.OrderCancelReq) response.Response[*data.OrderCancelRes]
}

type orderSvc struct {
	pg       db.PostgresqlDb
	repo     repository.Repo
	redisDB  redis_db.RedisDB
	producer kafka.Producer
}

func NewOrderSvc() OrderSvc {
	return &orderSvc{
		pg:       db.NewPostgresqlDb(),
		repo:     repository.NewRepo(),
		redisDB:  redis_db.NewRedisDB(),
		producer: kafka.NewProducer(),
	}
}

func (svc *orderSvc) Create(ctx context.Context, req data.OrderCreateReq) response.Response[*data.OrderCreateRes] {
	if valMap := validation.ValidateRequest(req); len(valMap) > 0 {
		return response.ValidationFailed[*data.OrderCreateRes](valMap)
	}

	// check product
	productIDs := make([]int64, len(req.OrderDetails))
	for i, detail := range req.OrderDetails {
		productIDs[i] = detail.ProductID
	}
	products, err := svc.repo.ProductRepo.GetPriceByIDs(productIDs...)
	if err != nil {
		return response.InternalServerError[*data.OrderCreateRes](err, "error get product")
	}

	// check match product
	orderDetails := make([]model.OrderDetail, len(products))
	for i := range products {
		if products[i].ID.IsNull() {
			return response.Notfound[*data.OrderCreateRes]("product not found")
		}
		orderDetails[i] = req.OrderDetails[i].ToOrderDetailDB(products[i])
		req.Total += orderDetails[i].SubTotal
	}

	var res data.OrderCreateRes

	if req.ID == nil {
		// insert
		err = svc.pg.ExecTx(ctx, func(tx db.TX) error {
			// insert order
			orderID, err := svc.repo.OrderRepo.Insert(tx, req.ToInsertOrderDB())
			if err != nil {
				return err
			}

			// insert order detail
			countOrderDetail, err := svc.repo.OrderDetailRepo.InsertMany(tx, orderID, orderDetails)
			if err != nil {
				return err
			}

			res.CreatedID = orderID
			res.CreateCount = countOrderDetail
			return nil

		})
		if err != nil {
			return response.InternalServerError[*data.OrderCreateRes](err, "error create order")
		}

	} else {
		// update
		// check order id
		dataDB, err := svc.repo.OrderRepo.GetByID(*req.ID)
		if err != nil {
			return response.InternalServerError[*data.OrderCreateRes](err, "error get order")
		}
		if dataDB == nil {
			return response.Notfound[*data.OrderCreateRes]("order not found")
		}

		err = svc.pg.ExecTx(ctx, func(tx db.TX) error {
			// clear order detail
			res.DeleteCount, err = svc.repo.OrderDetailRepo.DeleteMany(tx, dataDB.ID)
			if err != nil {
				return nil
			}

			// insert order detail
			res.CreateCount, err = svc.repo.OrderDetailRepo.InsertMany(tx, dataDB.ID, orderDetails)
			if err != nil {
				return err
			}

			// update order
			res.UpdateCount, err = svc.repo.OrderRepo.Update(tx, req.ToUpdateOrderDB(*dataDB))
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return response.InternalServerError[*data.OrderCreateRes](err, "error update order")
		}

	}

	go func() {
		// clear cache
		svc.redisDB.DeletePPrefix(fmt.Sprintf("%s:%d:", "orderSvc", req.UserID))
	}()

	res.CreateBy = req.CreateBy
	res.CreatedDate = time.Now().UTC()
	return response.Ok(&res)
}

func (svc *orderSvc) GetHistory(ctx context.Context, req data.OrderGetHistoryReq) response.Response[*data.OrderGetHistoryRes] {
	req = req.CleanReq()
	if valMap := validation.ValidateRequest(req); len(valMap) > 0 {
		return response.ValidationFailed[*data.OrderGetHistoryRes](valMap)
	}
	if valMap := req.ToVal(); len(valMap) > 0 {
		return response.ValidationFailed[*data.OrderGetHistoryRes](valMap)
	}

	var res data.OrderGetHistoryRes

	initial := req.IsInitial()
	cacheKey := fmt.Sprintf("%s:%d:%s", "orderSvc", req.UserID, "getHistory")
	if initial {
		if svc.redisDB.GetInfo(cacheKey, &res) == nil {
			return response.Ok(&res)
		}
	}

	dataDB, total, err := svc.repo.OrderRepo.GetHistory(req)
	if err != nil {
		return response.InternalServerError[*data.OrderGetHistoryRes](err, "error get order history")
	}

	checkResult := make(map[int64]data.OrderGetHistoryResult)
	var orderIDs []int64
	for i := range dataDB {
		if value, ok := checkResult[dataDB[i].ID]; !ok {
			orderIDs = append(orderIDs, dataDB[i].ID)
			checkResult[dataDB[i].ID] = data.OrderGetHistoryResult{}.FromDB(dataDB[i])
		} else {
			value.OrderDetails = append(value.OrderDetails, data.OrderDetailResult{}.FromDB(dataDB[i]))
			checkResult[dataDB[i].ID] = value
		}
	}

	results := make([]data.OrderGetHistoryResult, len(orderIDs))
	for i := range orderIDs {
		results[i] = checkResult[orderIDs[i]]
	}

	res = data.OrderGetHistoryRes{
		Results:      results,
		TotalPages:   int64(math.Ceil(float64(total) / float64(req.PerPage))),
		TotalResults: total,
		Page:         req.Page,
		PerPage:      req.PerPage,
	}

	if len(dataDB) > 0 {
		if initial {
			svc.redisDB.Set(cacheKey, res, 5*time.Minute)
		}
	}

	return response.Ok(&res)
}

func (svc *orderSvc) Delete(ctx context.Context, req data.OrderDeleteReq) response.Response[*data.OrderDeleteRes] {
	var res data.OrderDeleteRes

	err := svc.pg.ExecTx(ctx, func(tx db.TX) error {
		countDelete, err := svc.repo.OrderRepo.Delete(tx, req.ID)
		if err != nil {
			return err
		}
		res.DeletedCount = countDelete

		return nil
	})
	if err != nil {
		return response.InternalServerError[*data.OrderDeleteRes](err, "error delete order")
	}

	if res.DeletedCount > 0 {
		go func() {
			// clear cache
			svc.redisDB.DeletePPrefix(fmt.Sprintf("%s:%d:", "orderSvc", req.UserID))
		}()
	}

	res.DeletedBy = req.DeletedBy
	res.DeletedDate = time.Now().UTC()

	return response.Ok(&res)
}

func (svc *orderSvc) Confirm(ctx context.Context, req data.OrderConfirmReq) response.Response[*data.OrderConfirmRes] {

	dataDB, err := svc.repo.OrderRepo.GetByID(req.ID)
	if err != nil {
		return response.InternalServerError[*data.OrderConfirmRes](err, "error get order")
	}
	if dataDB == nil {
		return response.Notfound[*data.OrderConfirmRes]("order not found")
	}
	if dataDB.Status != meta.ORDER_STATUS_EDITABLE {
		return response.BadRequest[*data.OrderConfirmRes]("order not editable")
	}

	go func() {
		svc.producer.Send(events.OrderEvent{
			Role:       "admin",
			SendID:     1,
			OrderID:    dataDB.ID,
			Status:     meta.ORDER_STATUS_CONFIRM,
			UserID:     dataDB.UserID,
			Total:      dataDB.Total,
			ModifiedBy: req.ModifiedBy,
		})
	}()

	res := data.OrderConfirmRes{
		ModifiedBy:   req.ModifiedBy,
		ModifiedDate: time.Now().UTC(),
	}

	return response.Ok(&res)

}
func (svc *orderSvc) Shipping(ctx context.Context, req data.OrderShippingReq) response.Response[*data.OrderShippingRes] {
	dataDB, err := svc.repo.OrderRepo.GetByID(req.ID)
	if err != nil {
		return response.InternalServerError[*data.OrderShippingRes](err, "error get order")
	}
	if dataDB == nil {
		return response.Notfound[*data.OrderShippingRes]("order not found")
	}
	if dataDB.Status != meta.ORDER_STATUS_CONFIRM {
		return response.BadRequest[*data.OrderShippingRes]("order not editable")
	}

	go func() {
		svc.producer.Send(events.OrderEvent{
			Role:       "user",
			SendID:     dataDB.UserID,
			OrderID:    dataDB.ID,
			Status:     meta.ORDER_STATUS_SHIPPING,
			UserID:     dataDB.UserID,
			Total:      dataDB.Total,
			ModifiedBy: req.ModifiedBy,
		})
	}()

	res := data.OrderShippingRes{
		ModifiedBy:   req.ModifiedBy,
		ModifiedDate: time.Now().UTC(),
	}

	return response.Ok(&res)
}
func (svc *orderSvc) Completed(ctx context.Context, req data.OrderCompletedReq) response.Response[*data.OrderCompletedRes] {
	dataDB, err := svc.repo.OrderRepo.GetByID(req.ID)
	if err != nil {
		return response.InternalServerError[*data.OrderCompletedRes](err, "error get order")
	}
	if dataDB == nil {
		return response.Notfound[*data.OrderCompletedRes]("order not found")
	}
	if dataDB.Status != meta.ORDER_STATUS_SHIPPING {
		return response.BadRequest[*data.OrderCompletedRes]("order not editable")
	}

	go func() {
		svc.producer.Send(events.OrderEvent{
			Role:       "user",
			SendID:     dataDB.UserID,
			OrderID:    dataDB.ID,
			Status:     meta.ORDER_STATUS_COMPLETED,
			UserID:     dataDB.UserID,
			Total:      dataDB.Total,
			ModifiedBy: req.ModifiedBy,
		})
	}()

	res := data.OrderCompletedRes{
		ModifiedBy:   req.ModifiedBy,
		ModifiedDate: time.Now().UTC(),
	}

	return response.Ok(&res)
}
func (svc *orderSvc) Cancel(ctx context.Context, req data.OrderCancelReq) response.Response[*data.OrderCancelRes] {
	req.Reason = strings.TrimSpace(req.Reason)

	dataDB, err := svc.repo.OrderRepo.GetByID(req.ID)
	if err != nil {
		return response.InternalServerError[*data.OrderCancelRes](err, "error get order")
	}
	if dataDB == nil {
		return response.Notfound[*data.OrderCancelRes]("order not found")
	}
	if dataDB.Status != meta.ORDER_STATUS_SHIPPING {
		return response.BadRequest[*data.OrderCancelRes]("order not editable")
	}

	go func() {
		svc.producer.Send(events.OrderEvent{
			Role:       "user",
			SendID:     dataDB.UserID,
			OrderID:    dataDB.ID,
			Status:     meta.ORDER_STATUS_CANCEL,
			UserID:     dataDB.UserID,
			Reason:     req.Reason,
			Total:      dataDB.Total,
			ModifiedBy: req.ModifiedBy,
		})
	}()

	res := data.OrderCancelRes{
		ModifiedBy:   req.ModifiedBy,
		ModifiedDate: time.Now().UTC(),
	}

	return response.Ok(&res)

}
