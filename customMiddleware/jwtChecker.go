package customMiddleware

import (
	"Bancho/auth"
	"Bancho/userDB"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"
)
const rediskey = `config.jwt.thftgr` // 15



func JwtChecker() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			a := c.Request().Header.Get("Authorization")

			claims := auth.GetClaim(&a)
			if claims == nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			mySigningKey := []byte(userDB.Redis[15].HGet(context.TODO(), rediskey, "key").Val())
			t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.MapClaim).SignedString(mySigningKey)
			if err != nil {
				log.Error(err)

				return c.NoContent( http.StatusUnauthorized )
			}
			if c.Request().Header.Get("Authorization") != "Bearer "+t {
				return c.NoContent( http.StatusUnauthorized )
			}

			exp := userDB.Redis[1].HGet(context.TODO(), "account.token:"+claims.Jti, "exp").Val()
			if i, err := strconv.ParseInt(exp,10 ,64); err != nil || i < time.Now().Unix() {
				return c.NoContent( http.StatusUnauthorized )
			}

			return next(c)
		}
	}
}

