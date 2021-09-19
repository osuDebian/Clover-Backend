package mailHandler

import (
	"Bancho/userDB"
	"context"
	"fmt"
	"net/smtp"
)
const verifyBody = "From: No-replay\r\n"+
	"To: %s\r\n"+
	"Subject: MailVerifyCode\r\n%s"


func SendMailVerify(address, code *string)(err error) {
	ctx := context.Background()
	redis := userDB.Redis

	Id       := redis[15].HGet(ctx, userDB.KeyMail, "id").Val()
	//Sender   := redis[15].HGet(ctx, userDB.KeyMail, "sender").Val()
	Password := redis[15].HGet(ctx, userDB.KeyMail, "password").Val()
	Url      := redis[15].HGet(ctx, userDB.KeyMail, "url").Val()
	Tls      := redis[15].HGet(ctx, userDB.KeyMail, "tls").Val()
	//Ssl      := redis[15].HGet(ctx, userDB.KeyMail, "ssl").Val()


	auth := smtp.PlainAuth("", Id, Password, Url)
	body:= fmt.Sprintf(verifyBody,*address,*code)
	return smtp.SendMail(Tls, auth, Id, []string{*address}, []byte(body))
}
