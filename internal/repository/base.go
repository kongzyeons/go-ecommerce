package repository

import (
	"app-ecommerce/pkg/db"
	"sync"
)

type Repo struct {
	UserRepo        UserRepo
	ProductRepo     ProductRepo
	OrderRepo       OrderRepo
	OrderDetailRepo OrderDetailRepo
}

var repoInstance Repo
var repoOnce sync.Once

func NewRepo() Repo {
	repoOnce.Do(func() {
		pg := db.NewPostgresqlDb()
		repoInstance = Repo{
			UserRepo:        NewUserRepo(pg),
			ProductRepo:     NewProductRepo(pg),
			OrderRepo:       NewOrderRepo(pg),
			OrderDetailRepo: NewOrderDetailRepo(pg),
		}
	})
	return repoInstance
}
