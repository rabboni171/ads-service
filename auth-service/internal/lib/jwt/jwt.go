package jwt

import (
	"auth-service/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Создаем новый JWT токен, используя данные пользователя и конкретного приложения
func NewToken(user *models.User, app *models.App, duration time.Duration) (string, error) {  
    token := jwt.New(jwt.SigningMethodHS256)  

    // Добавляем в токен всю необходимую информацию
    claims := token.Claims.(jwt.MapClaims)  
    claims["uid"] = user.ID  
    claims["email"] = user.Email  
    claims["exp"] = time.Now().Add(duration).Unix()  
    claims["app_id"] = app.ID  

    // Подписываем токен, используя секретный ключ приложения
    tokenString, err := token.SignedString([]byte(app.Secret))  
    if err != nil {  
       return "", err  
    }  

    return tokenString, nil  
}