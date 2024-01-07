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

type CustomClaims struct {
	jwt.StandardClaims
	UserID   string
	ClientIP string
}

type Jwt struct {
	SigningKey []byte
}

func NewJWT() *Jwt {
	return &Jwt{SigningKey: []byte(cfg.App.JWTSecret)}
}

func (j *Jwt) GenerateToken(userID string, clientIP string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
		UserID:   userID,
		ClientIP: clientIP,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(cfg.App.JWTSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *Jwt) ParseToken(tokenStr string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(tokenStr *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if result, ok := err.(jwt.ValidationError); ok {
			switch {
			case result.Errors&jwt.ValidationErrorMalformed != 0:
				return EmptyString, EmptyString, errors.New(TokenMalformed)
			case result.Errors&jwt.ValidationErrorExpired != 0:
				return EmptyString, EmptyString, errors.New(TokenExpired)
			case result.Errors&jwt.ValidationErrorNotValidYet != 0:
				return EmptyString, EmptyString, errors.New(TokenNotValidYet)
			default:
				return EmptyString, EmptyString, errors.New(TokenInvalid)
			}
		}
		return EmptyString, EmptyString, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, claims.ClientIP, nil
	}
	return EmptyString, EmptyString, errors.New(TokenInvalid)
}

func (j *Jwt) RefreshToken(tokenStr string) (string, error) {
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
		return j.GenerateToken(claims.UserID, claims.ClientIP, duration)
	}
	return "", errors.New(TokenInvalid)
}
