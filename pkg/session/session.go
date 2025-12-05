package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

type SessionStorage interface {
	Get(c *fiber.Ctx, res interface{}) error
	Set(c *fiber.Ctx, key interface{}, value interface{}) error
	Delete(c *fiber.Ctx) error
}

type sessionStorage struct {
	sesionStore *session.Store
}

var sessionStorageInstance SessionStorage
var sessionStorageOnce sync.Once

func NewSessionStorage() SessionStorage {
	sessionStorageOnce.Do(func() {
		sessionStorageInstance = &sessionStorage{
			sesionStore: sesionStore,
		}
	})
	return sessionStorageInstance
}

func (s *sessionStorage) Get(c *fiber.Ctx, res interface{}) error {
	if res == nil {
		return errors.New("value invalid")
	}

	sess, err := s.sesionStore.Get(c)
	if err != nil {
		return err
	}

	userInfo := sess.Get("userInfo")
	if userInfo == nil {
		return errors.New("user info not found")
	}

	return json.Unmarshal([]byte(userInfo.(string)), &res)
}

func (s *sessionStorage) Set(c *fiber.Ctx, key interface{}, value interface{}) error {
	if value == nil || key == nil {
		return errors.New("value invalid")
	}

	// generate custom key
	s.sesionStore.KeyGenerator = func() string {
		return fmt.Sprintf("%s:%v:%s", "authSvc", key, utils.UUIDv4())
	}

	sess, err := s.sesionStore.Get(c)
	if err != nil {
		return err
	}

	newValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	sess.Set("userInfo", string(newValue))
	return sess.Save()
}

func (s *sessionStorage) Delete(c *fiber.Ctx) error {
	sess, err := s.sesionStore.Get(c)
	if err != nil {
		return err
	}
	return sess.Destroy()
}
