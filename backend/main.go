package main

import (
	"cmp"
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h25s_05/backend/handler"
	"github.com/traP-jp/h25s_05/backend/llm"
	"github.com/traP-jp/h25s_05/backend/llm/mock"
	"github.com/traP-jp/h25s_05/backend/llm/openai"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)
	conf := mysql.Config{
		User:                 cmp.Or(os.Getenv("DB_USER"), "root"),
		Passwd:               cmp.Or(os.Getenv("DB_PASSWORD"), "password"),
		Net:                  "tcp",
		Addr:                 cmp.Or(os.Getenv("DB_HOST"), "127.0.0.1") + ":" + cmp.Or(os.Getenv("DB_PORT"), "3306"),
		DBName:               cmp.Or(os.Getenv("DB_NAME"), "hackathon"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err := sqlx.Connect("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	llmsvc := provideLLMService()

	e := echo.New()
	config := handler.Config{
		DefaultAIIconURL: os.Getenv("DEFAULT_AI_ICON_URL"),
	}
	h := handler.NewHandler(config, db, llmsvc)
	h.SetUpRoutes(e.Group("/api"))
	e.Logger.Fatal(e.Start(":" + cmp.Or(os.Getenv("PORT"), "8080")))
}

func provideLLMService() llm.Service {
	if os.Getenv("DEBUG") == "1" {
		return mock.NewService()
	}
	mcp := llm.MCP{
		ServerLabel: "hackathon",
		ServerURL:   os.Getenv("MCP_SERVER_URL"),
		Header: map[string]string{
			"Authorization": os.Getenv("MCP_SERVER_AUTHORIZATION"),
		},
	}
	return openai.NewService([]llm.MCP{mcp})
}
