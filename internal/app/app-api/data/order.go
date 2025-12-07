package data

import (
	"app-ecommerce/internal/meta"
	"app-ecommerce/internal/model"
	"app-ecommerce/pkg/types"
	"reflect"
	"strings"
	"time"
)

type OrderDetail struct {
	ProductID int64 `json:"productID" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"required"`
}

type OrderCreateReq struct {
	ID           *int64        `json:"id"`
	Total        float64       `json:"-"`
	OrderDetails []OrderDetail `json:"orderDetails" validate:"required,max=10"`
	UserID       int64         `json:"-" validate:"required"`
	CreateBy     string        `json:"-" validate:"required"`
}

func (obj *OrderCreateReq) ToInsertOrderDB() model.Order {
	return model.Order{
		UserID:       obj.UserID,
		Total:        obj.Total,
		Status:       meta.ORDER_STATUS_EDITABLE,
		CreatedBy:    obj.CreateBy,
		CreatedDate:  time.Now(),
		ModifiedBy:   obj.CreateBy,
		ModifiedDate: types.NewNullTime(time.Now().UTC()),
	}
}

func (obj *OrderCreateReq) ToUpdateOrderDB(dataDB model.Order) model.Order {
	dataDB.Total = obj.Total
	dataDB.ModifiedBy = obj.CreateBy
	dataDB.ModifiedDate = types.NewNullTime(time.Now().UTC())
	return dataDB
}

func (obj *OrderDetail) ToOrderDetailDB(products model.GetPriceByIDsRes) model.OrderDetail {
	return model.OrderDetail{
		ProductID: products.ID.Int64,
		Quantity:  obj.Quantity,
		Price:     products.Price.Float64,
		SubTotal:  products.Price.Float64 * float64(obj.Quantity),
	}
}

type OrderCreateRes struct {
	CreatedID   int64     `json:"createdID"`
	CreateCount int64     `json:"createCount"`
	UpdateCount int64     `json:"updateCount"`
	DeleteCount int64     `json:"deleteCount"`
	CreateBy    string    `json:"createBy"`
	CreatedDate time.Time `json:"createdDate"`
}

type OrderGetHistoryReq struct {
	UserID  int64  `json:"userID" example:"0"`
	Status  string `json:"status" example:""`
	Page    int64  `json:"page" example:"1" validate:"gte=1"`
	PerPage int64  `json:"perPage" example:"10" validate:"gte=1"`
	SortBy  struct {
		Field     string       `json:"field" example:"modifiedDate"`
		FieldType reflect.Kind `json:"-"`
		Mode      string       `json:"mode" example:"desc"`
	} `json:"sortBy"`
}

func (obj *OrderGetHistoryReq) CleanReq() OrderGetHistoryReq {
	obj.Status = strings.TrimSpace(obj.Status)
	if obj.SortBy.Field == "" {
		obj.SortBy.Field = "modifiedDate"
		obj.SortBy.Mode = "desc"
	}
	return *obj
}

func (obj *OrderGetHistoryReq) ToVal() map[string]string {
	valMap := make(map[string]string)

	mapSortMode := map[string]bool{"ASC": true, "DESC": true}
	if _, ok := mapSortMode[strings.ToUpper(obj.SortBy.Mode)]; !ok {
		valMap["mode"] = "mode must be asc or desc"
	}

	if obj.Status != "" {
		mapOrderStatus := meta.GetOrderStatus()
		if _, ok := mapOrderStatus[obj.Status]; !ok {
			valMap["status"] = "status not found"
		}
	}

	type field struct {
		Field     string
		FieldType reflect.Kind
	}
	mapSort := map[string]field{
		"id":           {"o.id", reflect.Int64},
		"userID":       {"o.user_id", reflect.Int64},
		"total":        {"o.total", reflect.Int64},
		"status":       {"o.status", reflect.String},
		"modifiedDate": {"o.modified_date", reflect.Int64},
	}

	if value, ok := mapSort[obj.SortBy.Field]; ok {
		obj.SortBy.Field = value.Field
		obj.SortBy.FieldType = value.FieldType
	} else {
		valMap["field"] = "field not found"
	}

	return valMap
}

func (obj OrderGetHistoryReq) IsInitial() bool {
	if obj.Status == "" &&
		obj.Page == 1 &&
		obj.PerPage == 10 &&
		obj.SortBy.Field == "o.modified_date" &&
		obj.SortBy.Mode == "desc" {
		return true
	}
	return false
}

type OrderDetailResult struct {
	ID          int64   `json:"id"`
	ProudctName string  `json:"proudctName"`
	Description string  `json:"description"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

func (obj OrderDetailResult) FromDB(dataDB model.OrderGetInfoRes) OrderDetailResult {
	return OrderDetailResult{
		ID:          dataDB.OrderDetailID.Int64,
		ProudctName: dataDB.ProductName.String,
		Description: dataDB.Description.String,
		Quantity:    dataDB.Quantity.Int64,
		Price:       dataDB.Price.Float64,
		Subtotal:    dataDB.SubTotal.Float64,
	}
}

type OrderGetHistoryResult struct {
	ID           int64               `json:"id"`
	UserID       int64               `json:"userID"`
	Total        float64             `json:"total"`
	Status       string              `json:"status"`
	Reason       string              `json:"reason"`
	ModifiedDate *time.Time          `json:"modifiedDate"`
	OrderDetails []OrderDetailResult `json:"orderDetails"`
}

func (obj OrderGetHistoryResult) FromDB(dataDB model.OrderGetInfoRes) OrderGetHistoryResult {
	return OrderGetHistoryResult{
		ID:     dataDB.ID,
		UserID: dataDB.UserID,
		Total:  dataDB.Total,
		Status: dataDB.Status,
		Reason: dataDB.Reason.String,
		ModifiedDate: func() *time.Time {
			if dataDB.ModifiedDate.IsNull() {
				return nil
			}
			modifiedDate := dataDB.ModifiedDate.Time.UTC()
			return &modifiedDate
		}(),
		OrderDetails: []OrderDetailResult{{
			ID:          dataDB.OrderDetailID.Int64,
			ProudctName: dataDB.ProductName.String,
			Description: dataDB.Description.String,
			Quantity:    dataDB.Quantity.Int64,
			Price:       dataDB.Price.Float64,
			Subtotal:    dataDB.SubTotal.Float64,
		}},
	}
}

type OrderGetHistoryRes struct {
	Results      []OrderGetHistoryResult `json:"results"`
	TotalResults int64                   `json:"totalResults"`
	TotalPages   int64                   `json:"totalPages"`
	Page         int64                   `json:"page"`
	PerPage      int64                   `json:"perPage"`
}

type OrderDeleteReq struct {
	ID        int64  `json:"id" validate:"required"`
	UserID    int64  `json:"-" validate:"required"`
	DeletedBy string `json:"-" validate:"required"`
}

type OrderDeleteRes struct {
	DeletedCount int64     `json:"deletedCount"`
	DeletedBy    string    `json:"deletedBy"`
	DeletedDate  time.Time `json:"deletedDate"`
}

type OrderConfirmReq struct {
	ID int64 `json:"id" validate:"required"`

	UserID     int64  `json:"-" validate:"required"`
	ModifiedBy string `json:"-" validate:"required"`
}

type OrderConfirmRes struct {
	ModifiedBy   string    `json:"modifiedBy"`
	ModifiedDate time.Time `json:"modifiedDate"`
}

type OrderShippingReq struct {
	ID         int64  `json:"id" validate:"required"`
	ModifiedBy string `json:"-" validate:"required"`
}

type OrderShippingRes struct {
	ModifiedBy   string    `json:"modifiedBy"`
	ModifiedDate time.Time `json:"modifiedDate"`
}

type OrderCompletedReq struct {
	ID         int64  `json:"id" validate:"required"`
	ModifiedBy string `json:"-" validate:"required"`
}

type OrderCompletedRes struct {
	ModifiedBy   string    `json:"modifiedBy"`
	ModifiedDate time.Time `json:"modifiedDate"`
}

type OrderCancelReq struct {
	ID         int64  `json:"-" validate:"required"`
	Reason     string `json:"reason" validate:"required,max=255"`
	ModifiedBy string `json:"-" validate:"required"`
}

type OrderCancelRes struct {
	ModifiedBy   string    `json:"modifiedBy"`
	ModifiedDate time.Time `json:"modifiedDate"`
}
