package main

import (
	"cmp"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h25s_05/backend/handler"
)

func main() {
	conf := mysql.Config{
		User:                 cmp.Or(os.Getenv("DB_USER"), "root"),
		Passwd:               cmp.Or(os.Getenv("DB_PASSWORD"), "password"),
		Net:                  "tcp",
		Addr:                 cmp.Or(os.Getenv("DB_HOST"), "127.0.0.1") + ":" + cmp.Or(os.Getenv("DB_PORT"), "3306"),
		DBName:               cmp.Or(os.Getenv("DB_NAME"), "hackathon"),
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
	}
	db, err := sqlx.Connect("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	e := echo.New()
	h := handler.NewHandler(db)
	h.SetUpRoutes(e.Group("/api"))
	e.Logger.Fatal(e.Start(":" + cmp.Or(os.Getenv("PORT"), "8080")))
}


