package core

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	TokenExpired     = "Token已过期"
	TokenNotValidYet = "Token不再有效"
	TokenMalformed   = "Token非法"
	TokenInvalid     = "Token无效"
)

// JWT 从vault读取
//const JWTSecret = "123456"

type CustomClaims struct {
	jwt.StandardClaims
	Mobile string
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(cfg.App.JWTSecret)}
}

func (j *JWT) GenerateToken(text string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
		Mobile: text,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(cfg.App.JWTSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *JWT) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(tokenStr *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if result, ok := err.(jwt.ValidationError); ok {
			if result.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", errors.New(TokenMalformed)
			} else if result.Errors&jwt.ValidationErrorExpired != 0 {
				return "", errors.New(TokenExpired)
			} else if result.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", errors.New(TokenNotValidYet)
			} else {
				return "", errors.New(TokenInvalid)
			}
		}
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.Mobile, nil
	}
	return "", errors.New(TokenInvalid)
}

func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	duration := time.Hour * 24 * 7 // 一周
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(duration).Unix()
		return j.GenerateToken(claims.Mobile, duration)
	}
	return "", errors.New(TokenInvalid)
}
