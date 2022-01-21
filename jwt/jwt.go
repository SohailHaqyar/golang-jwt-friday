package jwt

import (
	"time"

	"github.com/SohailHaqyar/friday/data"
	"github.com/golang-jwt/jwt/v4"
)

func CreateJWTToken(user data.User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp

	t, err := token.SignedString([]byte("secret"))

	if err != nil {
		return "", 0, err
	}

	return t, exp, nil

}
