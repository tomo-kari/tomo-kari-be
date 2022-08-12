package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
	"tomokari/config"
	"tomokari/internal/constants"
)

type jwtClaims struct {
	jwt.StandardClaims
	Data interface{}
}

func GenerateToken(data interface{}, tokenType string) (signedToken string, err error) {
	var expiresAt int64
	now := time.Now().Local()
	nowUnix := now.Unix()
	fmt.Printf("%+v\n", nowUnix)
	cfg, err := config.NewConfig()
	if tokenType == constants.AccessToken {
		expiresAt = now.Add(time.Hour * time.Duration(cfg.JWT.ExpirationHours)).Unix()
	} else {
		expiresAt = now.Add(time.Hour * time.Duration(cfg.JWT.ExpirationHours*24*7)).Unix()
	}
	claims := &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "tomokari",
		},
		Data: data,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(cfg.JWT.SecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(signedToken string) (claims *jwtClaims, err error) {
	cfg, _ := config.NewConfig()
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)

	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil
}
