package jwt

import (
	"gin-derived/global"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserID   int64
	UserName string
	jwt.StandardClaims
}

func GetJwtSecret() []byte {
	return []byte(global.GCONFIG.Jwt.Secret)
}

func GenerateToken(userId int64, userName string) (string, error) {
	expire := global.GCONFIG.Jwt.Expire
	issuer := global.GCONFIG.Jwt.Issuer
	nowTime := time.Now()
	expireTime := nowTime.Add(expire)
	claims := Claims{
		UserID:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJwtSecret())
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
