package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	var data struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		writeJson(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	if data.Password != pass {
		writeJson(w, http.StatusOK, map[string]string{"error": "Неверный пароль"})
		return
	}

	hash := sha256.Sum256([]byte(pass))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": fmt.Sprintf("%x", hash),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte(pass))
	if err != nil {
		writeJson(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{"token": signedToken})
}
