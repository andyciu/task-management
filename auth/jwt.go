package auth

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pc01pc013/task-management/utils"
)

type JWTCustomClaims struct {
	Username string `json:"username"`
	AuthType uint   `json:"authtype"`
	jwt.StandardClaims
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Header.Authorization field.
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseResultFailed("Authorization is null in Header."))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseResultFailed("Format of Authorization is wrong."))
			return
		}
		// parts[0] is Bearer, parts is token.
		token, err := ParseToken(parts[1])

		if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
			log.Printf("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)
			c.Set("username", claims.Username)
			c.Set("authtype", claims.AuthType)
			c.Next()
		} else {
			log.Printf("ParseToken Error: %q", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MakeResponseResultFailed("Invalid Token."))
			return
		}
	}
}

func NewToken(username string, authtype uint) (string, error) {
	claims := JWTCustomClaims{
		username,
		authtype,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "AndyCiu",
		},
	}

	mySigningKey := []byte(os.Getenv("JWTSIGNKEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ParseToken(tokenstr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenstr, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWTSIGNKEY")), nil
	})
}
