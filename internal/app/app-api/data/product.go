package data

import (
	"app-ecommerce/internal/model"
	"reflect"
	"strings"
)

type ProductGetListReq struct {
	SearchText string `json:"searchText" example:"" validate:"max=255"`
	Page       int64  `json:"page" example:"1" validate:"gte=1"`
	PerPage    int64  `json:"perPage" example:"10" validate:"gte=1"`
	SortBy     struct {
		Field     string       `json:"field" example:"id"`
		FieldType reflect.Kind `json:"-"`
		Mode      string       `json:"mode" example:"asc"`
	} `json:"sortBy"`
}

func (obj *ProductGetListReq) CleanReq() ProductGetListReq {
	obj.SearchText = strings.TrimSpace(obj.SearchText)
	if obj.SortBy.Field == "" {
		obj.SortBy.Field = "id"
		obj.SortBy.Mode = "asc"
	}
	return *obj
}

func (obj *ProductGetListReq) ToVal() map[string]string {
	valMap := make(map[string]string)

	mapSortMode := map[string]bool{"ASC": true, "DESC": true}
	if _, ok := mapSortMode[strings.ToUpper(obj.SortBy.Mode)]; !ok {
		valMap["mode"] = "mode must be asc or desc"
	}

	type field struct {
		Field     string
		FieldType reflect.Kind
	}
	mapSort := map[string]field{
		"id":   {"id", reflect.Int64},
		"name": {"name", reflect.String},
	}

	if value, ok := mapSort[obj.SortBy.Field]; ok {
		obj.SortBy.Field = value.Field
		obj.SortBy.FieldType = value.FieldType
	} else {
		valMap["field"] = "field not found"
	}

	return valMap
}

func (obj ProductGetListReq) IsInitial() bool {
	if obj.SearchText == "" &&
		obj.Page == 1 &&
		obj.PerPage == 10 &&
		obj.SortBy.Field == "id" &&
		obj.SortBy.Mode == "asc" {
		return true
	}
	return false
}

type ProductGetListResult struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (obj ProductGetListResult) FromDB(dataDB model.Product) ProductGetListResult {
	return ProductGetListResult{
		ID:          dataDB.ID,
		Name:        dataDB.Name,
		Description: dataDB.Description,
		Price:       dataDB.Price,
	}
}

type ProductGetListRes struct {
	Results      []ProductGetListResult `json:"results"`
	TotalResults int64                  `json:"totalResults"`
	TotalPages   int64                  `json:"totalPages"`
	Page         int64                  `json:"page"`
	PerPage      int64                  `json:"perPage"`
}
