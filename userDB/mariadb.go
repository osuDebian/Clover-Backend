package userDB

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pterm/pterm"
)

var Maria *sql.DB

func ConnectMaria() {

	ctx := context.Background()
	id := Redis[15].HGet(ctx, KeySql, "user").Val()
	pw := Redis[15].HGet(ctx, KeySql, "password").Val()
	url:= Redis[15].HGet(ctx, KeySql, "address").Val()
	url += ":" + Redis[15].HGet(ctx, KeySql, "port").Val()

	db, err := sql.Open("mysql", id+":"+pw+"@tcp("+url+")/")
	if Maria = db; db != nil {
		Maria.SetMaxOpenConns(100)
		if _, err = Maria.Exec("SET SQL_SAFE_UPDATES = 0;"); err != nil {
			pterm.Error.Println("SET SQL_SAFE_UPDATES FAIL.", err)
			panic(err)
		}
		pterm.Info.Println("RDBMS Connected.")
	} else {
		pterm.Error.Println("RDBMS Connect Fail", err)
		panic(err)
	}
}
