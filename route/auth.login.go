package route

import "C"
import (
	"Bancho/auth"
	"Bancho/userDB"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequestBodyStruct struct {
	ID       string `json:"id" xml:"id" form:"id" `
	PASSWORD string `json:"password" xml:"password" form:"password" `
}



func Login(c echo.Context) (err error) {

	var req LoginRequestBodyStruct

	err = c.Bind(&req)
	if err != nil {
		_ = c.NoContent(http.StatusInternalServerError)
		return
	}

	var user LoginRequestBodyStruct
	err = userDB.Maria.QueryRow("SELECT ID,PASSWORD FROM BANCHO.ACCOUNT WHERE USER_LOGIN_ID = ? AND DELETED = FALSE", req.ID).Scan(&user.ID, &user.PASSWORD)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "UNKNOWN ACCOUNT OR DISABLED ACCOUNT. CHECK ID AND PASSWORD")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PASSWORD), []byte(req.PASSWORD)) != nil {
		return c.String(http.StatusBadRequest, "UNKNOWN ACCOUNT OR DISABLED ACCOUNT. CHECK ID AND PASSWORD")
	}
	//인증 성공
	body, err := auth.GenerateJwt(user.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "FAIL TO GEN TOKEN.")
	}

	return c.JSON(http.StatusOK, body)
}
