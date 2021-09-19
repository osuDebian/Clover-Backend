package userDB

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pterm/pterm"
	"os"
)

var Redis [16]*redis.Client

const (
	KeySql = `config.maria.debian`
	KeyMail = `config.mail.thftgr`
)

func RedisConnect() {
	for i := 0; i < 16; i++ {
		Redis[i] = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("X_REDIS_ADDRESS"),
			Username: os.Getenv("X_REDIS_USERNAME"),
			Password: os.Getenv("X_REDIS_PASSWORD"),
			DB:       i,
		})
	}

	if err := Redis[0].Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}

	pterm.Info.Println("CONNECTED REDIS.")

}
