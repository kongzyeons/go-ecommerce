package jwt

import (
	"app-ecommerce/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	UserID       int64         `json:"userID"`
	UserName     string        `json:"userName"`
	Email        string        `json:"email"`
	Role         string        `json:"role"`
	TimeDulation time.Duration `json:"timeDulation"`
}

func GenToken(req Token) (token string, err error) {
	cfg := config.GetConfig()
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   req.UserID,
		"userName": req.UserName,
		"email":    req.Email,
		"role":     req.Role,
		"exp":      time.Now().Add(req.TimeDulation).Unix(),
	},
	).SignedString([]byte(cfg.JwtKey))

	return token, err

}
