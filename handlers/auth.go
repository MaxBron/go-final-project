package handlers

import (
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

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtToken string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			var valid bool
			// здесь код для валидации и проверки JWT-токена
			todoPassword, _ := os.LookupEnv("TODO_PASSWORD")
			secret := []byte(todoPassword)
			jwtpassword := jwt.New(jwt.SigningMethodHS256)
			signedToken, _ := jwtpassword.SignedString(secret)
			if jwtToken == signedToken {
				valid = true
			}
			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

		}

		next(w, r)
	})

}
