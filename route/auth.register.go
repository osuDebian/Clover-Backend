package route

import (
	"Bancho/userDB"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"regexp"
)

var (
	idRegex *regexp.Regexp
)

func init() {
	idRegex, _ = regexp.Compile("^([A-Za-z0-9]{6,50})$")
}

type RegisterStruct struct {
	Id              string `json:"id"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	EmailVerifyCode string `json:"email_verify_code"`
}

func Register(c echo.Context) (err error) {
	var req RegisterStruct

	err = c.Bind(&req)
	if err != nil {
		return
	}
	if err = req.checkValue(); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err = userDB.Maria.Exec("INSERT INTO BANCHO.ACCOUNT (USER_LOGIN_ID, PASSWORD, USERNAME, EMAIL) VALUE (?,?,?,?)", req.Id, req.Password, req.Id, req.Email)
	if err != nil {
		_ = c.String(http.StatusInternalServerError, "RDBMS QUERY ERROR")
		return
	}
	pterm.Info.Println("NEW USER",req.Id,"JOIN")
	return c.NoContent(http.StatusOK )
}

func (v *RegisterStruct) checkValue() (err error) {

	if !idRegex.MatchString(v.Id) {
		return errors.New("ID NOT MATCH FORMAT")
	}
	if len(v.Password) < 8 {
		return errors.New("password to short")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(v.Password), 10)
	if err != nil {
		return err
	}
	v.Password = string(b)

	if _, err = mail.ParseAddress(v.Email); err != nil {
		return err
	}

	if v.EmailVerifyCode == "" || v.EmailVerifyCode != userDB.Redis[1].HGet(context.TODO(), redisKey+v.Email, "code").Val() {
		return errors.New("email Verify Fail")
	}
	pterm.Info.Println("email Verify successes.", userDB.Redis[1].Del(context.TODO(), redisKey+v.Email))

	return
}
