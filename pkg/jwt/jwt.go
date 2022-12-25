package jwt

import (
	"errors"
	"time"

	"github.com/abdullokh-mukhammadjonov/template_service/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func generate(claims jwt.MapClaims, exp int64, key interface{}) (string, error) {
	var token *jwt.Token
	token = jwt.New(jwt.SigningMethodHS256)

	claims["exp"] = exp

	token.Claims = claims

	return token.SignedString(key)
}

func GenerateJWT(userType, id string) (access, refresh string, err error) {
	cfg := config.Load()
	sessionID, _ := uuid.NewRandom()
	claims := make(jwt.MapClaims)

	claims["iss"] = "macbro"
	claims["session_id"] = sessionID
	claims["iat"] = time.Now().Unix()
	claims["user_type"] = userType
	claims["sub"] = id

	access, err = generate(claims, time.Now().Add(24*time.Hour).Unix(), cfg.SigningKey)

	if err != nil {
		err = errors.New("error while generating access_token")
		return
	}

	refresh, err = generate(claims, time.Now().AddDate(0, 0, 7).Unix(), cfg.RefreshSigningKey)

	if err != nil {
		err = errors.New("error while generating refresh_token")
		return
	}

	return
}

func ExtractClaims(tokenStr string, signingKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !(ok && token.Valid) {
		return nil, errors.New("invalid jwt token")
	}

	return claims, nil
}
