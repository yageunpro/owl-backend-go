package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/service"
	"github.com/yageunpro/owl-backend-go/service/auth"
	"net/http"
)

type Handler interface {
	GoogleLogin(c echo.Context) error
	GoogleCallback(c echo.Context) error
	DevSignUp(c echo.Context) error
	DevSignIn(c echo.Context) error
}
type handler struct {
	svc *service.Service
}

func New(svc *service.Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) GoogleLogin(c echo.Context) error {
	req := new(reqGoogleLogin)

	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}

	out, err := h.svc.Auth.GoogleLogin(c.Request().Context(), req.Ref)
	if err != nil {
		return err
	}

	c.SetCookie(out.Cookie)

	return c.Redirect(http.StatusTemporaryRedirect, out.RedirectURL)
}

func (h *handler) GoogleCallback(c echo.Context) error {
	req := new(reqGoogleCallback)
	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}

	cookie, err := c.Request().Cookie(auth.CookieKey)

	if err != nil {
		return err
	}

	out, err := h.svc.Auth.GoogleCallback(c.Request().Context(), auth.GoogleCallbackParam{
		Cookie: cookie,
		State:  req.State,
		Code:   req.Code,
		Scope:  req.Scope,
	})
	if err != nil {
		return err
	}

	_ = h.svc.Calendar.Sync(c.Request().Context(), out.UserId)

	c.SetCookie(out.Cookie)
	c.SetCookie(&http.Cookie{
		Name:   auth.CookieKey,
		MaxAge: -1,
	})

	return c.Redirect(http.StatusTemporaryRedirect, out.RedirectURL)
}

func (h *handler) DevSignUp(c echo.Context) error {
	req := new(reqDevSignUp)
	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}

	token, err := h.svc.Auth.DevSignUp(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	res := resDevSignUp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) DevSignIn(c echo.Context) error {
	req := new(reqDevSignIn)
	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}

	token, err := h.svc.Auth.DevSignIn(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	res := resDevSignIn{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return c.JSON(http.StatusOK, res)
}
