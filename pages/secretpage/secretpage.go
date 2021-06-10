package secretpage

import (
	"encoding/json"
	"fmt"
	"iitk-coin/model"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)


func AccessPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var result model.User
	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Name = claims["name"].(string)
		result.CanAccessPage=claims["accesspage"].(bool)
		if result.CanAccessPage {
			json.NewEncoder(w).Encode("Access allowed")
			return
		}
		json.NewEncoder(w).Encode("Access denied")
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}