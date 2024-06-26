package main

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/yageunpro/owl-backend-go/handler"
	"github.com/yageunpro/owl-backend-go/internal/config"
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/internal/google/oauth"
	"github.com/yageunpro/owl-backend-go/internal/jwt"
	"github.com/yageunpro/owl-backend-go/internal/naver/search"
	"github.com/yageunpro/owl-backend-go/service"
	"github.com/yageunpro/owl-backend-go/store"
	"log/slog"
	"os"
	"strings"
	"time"
)

func main() {
	pool, err := db.Connect(config.DBDsn)
	if err != nil {
		slog.Error("fail to init pool", "detail", err)
		os.Exit(1)
	}

	o, err := oauth.New(oauth.Config{
		ClientId:     config.OAuth.ClientId,
		ClientSecret: config.OAuth.ClientSecret,
		RedirectUri:  config.OAuth.RedirectUri,
		Scopes:       config.OAuth.Scopes,
	})
	if err != nil {
		slog.Error("fail to init oauth", "detail", err)
		os.Exit(1)
	}
	oauth.InitGlobal(o)

	n := search.New(config.Naver.ClientID, config.Naver.ClientSecret)
	search.InitGlobal(n)

	jwt.SetSecretKey(config.JWT.AccessKey, config.JWT.RefreshKey)

	sto, err := store.New(pool)
	if err != nil {
		slog.Error("fail to init store", "detail", err)
		os.Exit(1)
	}

	svc, err := service.New(sto)
	if err != nil {
		slog.Error("fail to init service", "detail", err)
		os.Exit(1)
	}

	hdl, err := handler.New(svc)
	if err != nil {
		slog.Error("fail to init handler", "detail", err)
		os.Exit(1)
	}

	schl, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("fail to init scheduler", "detail", err)
		os.Exit(1)
	}

	_, err = schl.NewJob(
		gocron.DurationJob(time.Hour),
		gocron.NewTask(
			func() {
				slog.Info("start sync calendars job")
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
				defer cancel()
				userIds, err := sto.Auth.GetAllOAuthUserIds(ctx)
				if err != nil {
					slog.Error("job failed, GetAllOAuthUserIds", "detail", err)
					return
				}

				for _, userId := range userIds {
					err = svc.Calendar.Sync(ctx, userId)
					if err != nil {
						slog.Error("job failed, Sync", "detail", err)
						return
					}
				}
				slog.Info("end sync calendars job")
				return
			}),
	)
	if err != nil {
		slog.Error("fail to init job", "detail", err)
		os.Exit(1)
	}

	schl.Start()
	defer schl.Shutdown()

	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Request().RequestURI, "/api/docs") {
				return true
			}
			return false
		},
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Debug = true
	e.HideBanner = true
	Route(e, hdl)
	if err := e.Start(":8000"); err != nil {
		panic(err)
	}
}
