package secretpage

import (
	"encoding/json"
	"fmt"
	"iitk-coin/model"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

func HasAccess(tokenString string, access_required string) (bool,error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if err!=nil {
		return false,err
	}
	var result model.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Name = claims["name"].(string)
		result.Access=claims["access"].(string)
		if result.Access==access_required {
			return true,nil
		}	
	}
	return false,nil
}
func StudentAccess(tokenString string) (int,error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method")
	}
	return []byte("secret"),nil
	})
	if err!=nil {
		return -1,err
	}
	var result model.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Name = claims["name"].(string)
		result.Rollno,err=strconv.Atoi(claims["rollno"].(string))
		result.Access=claims["access"].(string)
		fmt.Print(result.Rollno)
		return result.Rollno,nil
	}
	return -1,err
}
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
		result.Access=claims["access"].(string)
		if result.Access=="gensec" {
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