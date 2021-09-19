package main

import (
	"Bancho/Logger"
	"Bancho/customMiddleware"
	"Bancho/route"
	"Bancho/userDB"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pterm/pterm"
	"log"
	"net/http"
)

func init() {
	userDB.RedisConnect()
	userDB.ConnectMaria()

}

func main() {
	e := echo.New()
	e.Pre(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://thftgr.stoplight.io"},
		AllowMethods:     []string{"GET","POST","HEAD"},
	}))
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middleware.Logger())
	e.HideBanner = true
	go func() {
		for {
			<- Logger.Ch
			e.Logger.SetOutput(log.Writer())
			pterm.Info.Println("UPDATED ECHO LOGGER.")
		}
	}()



	e.GET("/jwt/test", func(c echo.Context) error {
		return c.String(http.StatusOK,"OK")
	},customMiddleware.JwtChecker())

	auth := e.Group("/auth")
	auth.POST("/login",route.Login)
	auth.POST("/register",route.Register)
	auth.POST("/logout",route.Logout)
	auth.POST("/verify/mail",route.EmailVerify)

	oauth := e.Group("/oauth")
	oauth.GET("/authorize",nil)
	oauth.POST("/token",nil)

	beatmaps := e.Group("/beatmaps")
	beatmaps.GET("/:beatmap/scores/users/:user",nil)
	beatmaps.GET("/:beatmap/scores",nil)
	beatmaps.GET("/:beatmap",nil)

	e.GET("/rankings/:mode/:type",route.Rankings)

	e.Logger.Fatal(e.Start(":80"))



}


