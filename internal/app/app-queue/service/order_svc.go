package service

import (
	"app-ecommerce/internal/events"
	"app-ecommerce/internal/repository"
	"app-ecommerce/pkg/db"
	"app-ecommerce/pkg/hub"
	redis_db "app-ecommerce/pkg/redis"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type orderSvc struct {
	pg      db.PostgresqlDb
	repo    repository.Repo
	webSse  hub.WebSse
	redisDB redis_db.RedisDB
}

func NewOrderSvc() events.EventHandler {
	return &orderSvc{
		pg:      db.NewPostgresqlDb(),
		repo:    repository.NewRepo(),
		webSse:  hub.NewWebSse(),
		redisDB: redis_db.NewRedisDB(),
	}
}

func (svc *orderSvc) Handle(topic string, eventBytes []byte) error {
	evnet := events.OrderEvent{}
	if topic != reflect.TypeOf(evnet).Name() {
		return errors.New("topic not match")
	}
	if err := json.Unmarshal(eventBytes, &evnet); err != nil {
		return err
	}

	err := svc.pg.ExecTx(context.Background(), func(tx db.TX) error {
		_, err := svc.repo.OrderRepo.Update(tx, evnet.ToOrderUpdateDB())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	jsonStr, err := json.Marshal(hub.Message{
		Topic:   topic,
		Role:    evnet.Role,
		SendID:  evnet.SendID,
		UserID:  evnet.UserID,
		Content: string(eventBytes),
	})
	if err != nil {
		return err
	}

	go func() {
		// clear cache
		svc.redisDB.DeletePPrefix(fmt.Sprintf("%s:%d:", "orderSvc", evnet.UserID))
	}()

	svc.webSse.Broadcast(string(jsonStr))

	return nil
}
