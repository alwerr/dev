package dev

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(`secret-key`)

type Claims struct {
	jwt.RegisteredClaims
	ID   string
	Name string
	nbf  int64
}

// Sign new token
func Sign(name, ids string) (string, string, error) {
	// var id = uuid.New()
	var id = ids
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			`Name`: name,
			`ID`:   id,
			`exp`:  time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return ``, ``, err
	}
	return tokenString, id, nil
	// return tokenString, id.String(), nil
}

// return the access token
func Signed(r *http.Request) (Claims, bool) {
	token, err := jwt.ParseWithClaims(r.Header.Get(`Authorization`), &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return Claims{}, true
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return Claims{ID: claims.ID, Name: claims.Name}, false
	}
	return Claims{}, true
}
