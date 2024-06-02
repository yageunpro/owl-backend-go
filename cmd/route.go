package main

import (
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/handler"
	"github.com/yageunpro/owl-backend-go/internal/jwt"
	"github.com/yageunpro/owl-backend-go/internal/openapi"
	"net/http"
)

func Route(e *echo.Echo, h *handler.Handler) {
	g := e.Group("/api", jwt.Middleware)

	g.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/api/docs/index.html")
	})
	g.GET("/docs/*", echo.WrapHandler(openapi.Handler()))

	auth := g.Group("/auth")

	appointment := g.Group("/appointment")
	appointment.Use()
	calendar := g.Group("/calendar")
	calendar.Use()
	location := g.Group("/location")
	location.Use()
	user := g.Group("/user")
	user.Use()

	auth.GET("/oauth/google", h.Auth.GoogleLogin)
	auth.GET("/callback/google", h.Auth.GoogleCallback)
	auth.POST("/dev/signup", h.Auth.DevSignUp)
	auth.POST("/dev/signin", h.Auth.DevSignIn)

	appointment.POST("", h.Appointment.Add)
	appointment.GET("/list", h.Appointment.List)
	appointment.GET("/:id", h.Appointment.Info)
	appointment.PATCH("/:id", h.Appointment.Edit)
	appointment.DELETE("/:id", h.Appointment.Delete)
	appointment.GET("/:id/share", h.Appointment.Share)
	appointment.POST("/:id/join", h.Appointment.Join)
	appointment.POST("/:id/join/nonmember", h.Appointment.JoinNonmember)
	appointment.GET("/:id/recommend", h.Appointment.RecommendTime)
	appointment.POST("/:id/confirm", h.Appointment.Confirm)

	calendar.POST("/schedule", h.Calendar.ScheduleAdd)
	calendar.GET("/schedule/:id", h.Calendar.ScheduleInfo)
	calendar.DELETE("/schedule/:id", h.Calendar.ScheduleDelete)
	calendar.GET("/schedule/list", h.Calendar.ScheduleList)
	calendar.POST("/sync", h.Calendar.Sync)

	location.GET("", h.Location.Search)

	user.GET("/me", h.User.Me)
	user.GET("/account", h.User.ListAccount)
	user.POST("/account", h.User.AddAccount)
	user.POST("/account/:id", h.User.VerifyAccount)
	user.DELETE("/account/:id", h.User.DeleteAccount)
}
