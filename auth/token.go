package auth

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var (
	// jwt 签名加密
	secret = viper.GetString("jwt.signedSecret")
)

// UserClaims UserClaims
type UserClaims struct {
	ID   int64
	Name string
	jwt.StandardClaims
}

// iss: 签发者
// sub: 面向的用户
// aud: 接收方
// exp: 过期时间
// nbf: 生效时间
// iat: 签发时间
// jti: 唯一身份标识

func generate(auth *UserClaims) (string, error) {
	auth.Id = strconv.FormatInt(auth.ID, 10)
	auth.Issuer = auth.Name
	auth.IssuedAt = time.Now().Unix()
	auth.ExpiresAt = time.Now().Add(time.Second * 10).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth)
	return token.SignedString(secret)
}

func parse(tokenStr string) (*UserClaims, error) {
	auth := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, auth, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Token is invalid, value: %s", tokenStr)
	}
	return auth, nil
}
