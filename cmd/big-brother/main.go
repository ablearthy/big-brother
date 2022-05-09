package main

import (
	"big-brother/internal/background"
	"big-brother/internal/config"
	"big-brother/internal/db"
	"big-brother/internal/postinit"
	"big-brother/internal/server"
	mytemplate "big-brother/internal/template"
	"context"
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func parseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yml", "path to config.yml")
	flag.Parse()
	return configPath, nil
}

func run() error {
	configPath, err := parseFlags()

	if err != nil {
		return err
	}

	cfg, err := config.NewConfig(configPath)

	if err != nil {
		return err
	}

	fmt.Println(cfg)

	err = db.Connect(cfg.Db)
	if err != nil {
		return err
	}

	defer db.GetConn().Close(context.Background())

	rand.Seed(time.Now().UnixNano())

	background.InitLongPollManagerWrapper()
	go background.GetLongPollManagerWrapper().Run()

	postinit.StartLongPollForAllUsers()

	tmpl := &mytemplate.Template{
		Templates: template.Must(template.ParseGlob("public/*.html")),
	}

	e := echo.New()
	e.Renderer = tmpl

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cfg.Server.CookiesSecret))))

	server.SetRoutes(e)

	e.Logger.Fatal(e.Start(":3000"))
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
