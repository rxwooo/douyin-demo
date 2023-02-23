package jwt

import (
	"errors"
	"fmt"
	"time"

	"crypto/md5"
	"encoding/hex"

	jwt "github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 2

// const TokenExpireDuration = time.Second * 60

var Secret = []byte("123456789")

type MyClaims struct {
	User_id int
	jwt.StandardClaims
}

// get token
func GetToken(user_id int) (string, error) {
	cla := MyClaims{
		user_id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "rxw-jwt",                                  // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)
	fmt.Println("Token = ", token)
	return token.SignedString(Secret) // 进行签名生成对应的token
}

// parse token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 中间件,认证token合法性
func JwtAuthMiddleware(token string) (int32, error) {

	//code:
	//10002 token valid
	//10003 token invalid

	_, err := ParseToken(token)
	if err != nil {
		fmt.Println("err = ", err.Error())
		return 10003, err
	}
	return 10002, nil
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
