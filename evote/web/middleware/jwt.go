package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var jwtkey = []byte("www.encode.com")

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Settoken(username string) string {
	myClaims := MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,
			ExpiresAt: time.Now().Unix() + 60*60*2,
			Issuer:    "evote",
			Subject:   "user token",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	token, err := t.SignedString(jwtkey)
	if err != nil {
		fmt.Println(err.Error())
	}
	return token
}

func Verfiy(tokenString string, ctx *gin.Context) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil || !token.Valid {
		ctx.Abort()
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "token expire", "code": "0"})
	} else {
		ctx.Next()
	}
}
