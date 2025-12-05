package repository

import (
	"app-ecommerce/pkg/db"
	"sync"
)

type Repo struct {
	UserRepo    UserRepo
	ProductRepo ProductRepo
}

var repoInstance Repo
var repoOnce sync.Once

func NewRpo() Repo {
	repoOnce.Do(func() {
		pg := db.NewPostgresqlDb()
		repoInstance = Repo{
			UserRepo:    NewUserRepo(pg),
			ProductRepo: NewProductRepo(pg),
		}
	})
	return repoInstance
}
