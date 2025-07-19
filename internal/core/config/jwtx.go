package config

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpired     = "Token已过期"
	TokenNotValidYet = "Token不再有效"
	TokenMalformed   = "Token非法"
	TokenInvalid     = "Token无效"
)

type CustomClaims struct {
	jwt.RegisteredClaims

	UserID   string
	ClientIP string
}

type Jwt struct {
	SigningKey []byte
}

func NewJWT() *Jwt {
	return &Jwt{SigningKey: []byte(cfg.App.JWTSecret)}
}

func (j *Jwt) GenerateToken(
	userID string,
	clientIP string,
	duration time.Duration,
) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now),               // 签发时间
			NotBefore: jwt.NewNumericDate(now),               // 生效时间
			Issuer:    cfg.App.Name,
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
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return "", "", errors.New(TokenMalformed)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", errors.New(TokenExpired)
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return "", "", errors.New(TokenNotValidYet)
		} else {
			return "", "", errors.New(TokenInvalid)
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, claims.ClientIP, nil
	}

	return "", "", errors.New(TokenInvalid)
}

func (j *Jwt) RefreshToken(tokenStr string) (string, error) {
	duration := time.Hour * 24 * 7 // 一周

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		},
		jwt.WithTimeFunc(time.Now), // 设置时间函数
	)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))

		return j.GenerateToken(claims.UserID, claims.ClientIP, duration)
	}

	return "", errors.New(TokenInvalid)
}
