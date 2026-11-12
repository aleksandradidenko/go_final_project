package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func auth(next http.HandlerFunc, pass string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		if len(pass) > 0 {
			var tokenString string

			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}

			var valid bool

			// проверяем JWT-токен
			if tokenString != "" {
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if token.Method != jwt.SigningMethodHS256 {
						return nil, fmt.Errorf("неверный метод подписи")
					}
					return []byte(pass), nil
				})

				if err == nil && token.Valid {
					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						hash := sha256.Sum256([]byte(pass))
						valid = claims["hash"] == fmt.Sprintf("%x", hash)
					}
				}
			}

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
