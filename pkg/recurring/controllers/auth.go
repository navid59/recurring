package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("mySuperSecrtePhrase")

func IsAuthorized(c *gin.Context) {

	token := c.Request.Header.Get("token")
	if len(token) > 0 {

		_, err := jwt.Parse(c.Request.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return mySigningKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": err.Error(),
			})
		}

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Missing token!",
		})
	}
}

func IsLicensed(tokenStr, MerchantSignature string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if _, ok := claims["aud"]; ok && claims["aud"] == MerchantSignature {
			return claims, true
		}

		return claims, false
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
