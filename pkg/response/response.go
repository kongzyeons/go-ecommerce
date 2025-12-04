package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Response[T any] struct {
	Success          bool              `json:"success"`
	StatusCode       int               `json:"statusCode"`
	Message          string            `json:"message"`
	ValidationErrors map[string]string `json:"validationErrors"`
	Data             T                 `json:"data"`
}

func (r Response[T]) JSON(c *fiber.Ctx) error {
	return c.Status(r.StatusCode).JSON(r)
}

func Ok[T any](data T, msgs ...string) Response[T] {
	if len(msgs) > 0 {
		jsonStr, err := json.Marshal(data)
		if err == nil {
			msgs[0] = fmt.Sprintf("%s : %s", msgs[0], string(jsonStr))
		}
		log.Println(msgs[0])
	}

	return Response[T]{
		Success:    true,
		StatusCode: 200,
		Data:       data,
	}
}

func BadRequest[T any](msgs ...string) Response[T] {
	message := http.StatusText(http.StatusBadRequest)
	if len(msgs) > 0 {
		message = strings.Join(msgs, " ")
	}
	log.Println(message)
	return Response[T]{
		Success:    false,
		StatusCode: 400,
		Message:    message,
	}
}

func ValidationFailed[T any](valMap map[string]string) Response[T] {
	message := "Validate Bad Request"
	jsonStr, err := json.Marshal(valMap)
	if err == nil {
		message = string(jsonStr)
	}
	log.Println(message)

	return Response[T]{
		Success:          false,
		StatusCode:       400,
		Message:          "Validate Bad Request",
		ValidationErrors: valMap,
	}
}

func Unauthorized[T any](msgs ...string) Response[T] {
	message := http.StatusText(http.StatusUnauthorized)
	if len(msgs) > 0 {
		message = strings.Join(msgs, " ")
	}
	log.Println(message)
	return Response[T]{
		Success:    false,
		StatusCode: 401,
		Message:    message,
	}
}

func Notfound[T any](message string) Response[T] {
	if message == "" {
		message = "Not Found"
	}
	log.Println(message)

	return Response[T]{
		Success:    false,
		StatusCode: 404,
		Message:    message,
	}
}

func InternalServerError[T any](err error, message string) Response[T] {
	log.Println(message, err)
	return Response[T]{
		Success:    false,
		StatusCode: 500,
		Message:    message,
	}
}
