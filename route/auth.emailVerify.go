package route

import (
	"Bancho/common"
	"Bancho/mailHandler"
	"Bancho/userDB"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"net/http"
	"net/mail"
	"time"
)

const redisKey = "Bancho.email.verify:"

type EmailVerifyStruct struct {
	Address string `json:"address"`
}

func EmailVerify(c echo.Context) (err error) {
	var req EmailVerifyStruct
	err = c.Bind(&req)
	if err != nil {
		_ = c.NoContent(http.StatusBadRequest)
		return
	}
	_, err = mail.ParseAddress(req.Address)
	if err != nil {
		_ = c.String(http.StatusBadRequest, err.Error())
		return
	}

	code := common.GenEmailCode() // 메일인증코드 생성
	ctx := context.TODO()

	err = userDB.Redis[1].HSet(ctx, redisKey+req.Address, "code",code).Err() // 레디스 저장
	if err != nil {
		pterm.Error.Println(err)
		_ = c.NoContent(http.StatusInternalServerError)
		return
	}

	err = userDB.Redis[1].Expire(ctx, redisKey+req.Address, time.Minute * 10).Err() // 타임아웃 지정
	if err != nil {
		pterm.Error.Println(err)
		_ = c.NoContent(http.StatusInternalServerError)
		return
	}
	err = mailHandler.SendMailVerify(&req.Address,&code)
	if err != nil {
		pterm.Error.Println(err)
		_ = c.NoContent(http.StatusInternalServerError)

		err = userDB.Redis[1].Del(ctx, redisKey+req.Address).Err() // 메일 발송 실패시 삭제
		if err != nil {
			pterm.Error.Println(err)
			_ = c.NoContent(http.StatusInternalServerError)
			return
		}
		return
	}
	return
}
