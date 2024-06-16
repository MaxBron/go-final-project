package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	buf.ReadFrom(r.Body)
	user := User{}
	json.Unmarshal(buf.Bytes(), &user)
	todoPassword, _ := os.LookupEnv("TODO_PASSWORD")
	if user.Password == todoPassword {
		secret := []byte(user.Password)
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		signedToken, _ := jwtToken.SignedString(secret)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, signedToken)))
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Неверный пароль"}`))
	}

}
