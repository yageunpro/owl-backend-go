package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/yageunpro/owl-backend-go/internal/google/oauth"
	"github.com/yageunpro/owl-backend-go/internal/google/user"
	"github.com/yageunpro/owl-backend-go/internal/jwt"
	"github.com/yageunpro/owl-backend-go/store"
	"github.com/yageunpro/owl-backend-go/store/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Service interface {
	GoogleLogin(ctx context.Context, ref string) (*resGoogleLogin, error)
	GoogleCallback(ctx context.Context, arg GoogleCallbackParam) (*resGoogleCallback, error)
	DevSignUp(ctx context.Context, email, password string) (*resToken, error)
	DevSignIn(ctx context.Context, email, password string) (*resToken, error)
}

type service struct {
	store *store.Store
}

func New(sto *store.Store) Service {
	return &service{store: sto}
}

func (s *service) GoogleLogin(ctx context.Context, ref string) (*resGoogleLogin, error) {
	if ref == "" {
		ref = "/"
	}

	val := refCookieValue{
		State: uuid.New(),
		Ref:   ref,
	}

	url, err := oauth.AuthCodeURL(val.State.String(), true)
	if err != nil {
		return nil, errors.Join(errors.New("failed to generate auth url"), err)
	}

	out, err := json.Marshal(val)
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal ref cookie"), err)
	}

	cookie := &http.Cookie{
		Name:     CookieKey,
		Value:    base64.URLEncoding.EncodeToString(out),
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute).UTC(),
		HttpOnly: true,
	}

	return &resGoogleLogin{
		RedirectURL: url,
		Cookie:      cookie,
	}, nil
}

func (s *service) GoogleCallback(ctx context.Context, arg GoogleCallbackParam) (*resGoogleCallback, error) {
	val := new(refCookieValue)
	raw, err := base64.URLEncoding.DecodeString(arg.Cookie.Value)
	if err != nil {
		return nil, errors.Join(errors.New("failed to decode cookie value"), err)
	}
	err = json.Unmarshal(raw, val)
	if err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal ref cookie"), err)
	}

	if val.State.String() != arg.State {
		return nil, errors.New("failed to validate state")
	}

	token, err := oauth.Token(ctx, arg.Code)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get oauth token"), err)
	}

	info, err := user.GetUserInfo(token.AccessToken)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get user info"), err)
	}

	var userId uuid.UUID
	out, err := s.store.Auth.GetOAuthUser(ctx, info.OpenId)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			userId = uuid.Must(uuid.NewV7())
			var refreshToken *string
			if token.RefreshToken == "" {
				refreshToken = nil
			} else {
				refreshToken = &token.RefreshToken
			}
			isAllow, err := oauth.IsAllowSync(arg.Scope)
			if err != nil {
				return nil, errors.Join(errors.New("failed to check oauth permission"), err)
			}

			err = s.store.Auth.CreateOAuthUser(ctx, auth.CreateOAuthUserParam{
				UserId:       userId,
				Email:        info.Email,
				UserName:     info.Name,
				OpenId:       info.OpenId,
				AccessToken:  token.AccessToken,
				RefreshToken: refreshToken,
				AllowSync:    isAllow,
				ValidUntil:   token.Expiry,
			})
			if err != nil {
				return nil, errors.Join(errors.New("failed to create new user"), err)
			}
		} else {
			return nil, errors.Join(errors.New("failed to get user"), err)
		}
	} else {
		userId = out.UserId
		var refreshToken *string
		if token.RefreshToken != "" {
			refreshToken = &token.RefreshToken
		}
		isAllow, err := oauth.IsAllowSync(arg.Scope)
		if err != nil {
			return nil, errors.Join(errors.New("failed to check oauth permission"), err)
		}
		err = s.store.Auth.UpdateOAuthUser(ctx, auth.UpdateOAuthUserParam{
			UserId:       userId,
			OpenId:       info.OpenId,
			AccessToken:  token.AccessToken,
			RefreshToken: refreshToken,
			AllowSync:    isAllow,
			ValidUntil:   token.Expiry.UTC(),
		})
		if err != nil {
			return nil, errors.Join(errors.New("failed to update user"), err)
		}
	}

	accessToken, err := jwt.NewAccessToken(userId)
	if err != nil {
		return nil, errors.Join(errors.New("could not create access token"), err)
	}
	refreshToken, err := jwt.NewRefreshToken(userId)
	if err != nil {
		return nil, errors.Join(errors.New("could not create refresh token"), err)
	}

	authCookie, err := jwt.ToCookie(accessToken, refreshToken)
	if err != nil {
		return nil, errors.Join(errors.New("failed to generate auth cookie"), err)
	}

	res := resGoogleCallback{
		RedirectURL: val.Ref,
		Cookie:      authCookie,
	}
	return &res, nil
}

func (s *service) DevSignUp(ctx context.Context, email, password string) (*resToken, error) {
	userId := uuid.Must(uuid.NewV7())
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Join(errors.New("could not hash password"), err)
	}

	err = s.store.Auth.CreateDevUser(ctx, auth.CreateDevUserParam{
		UserId:       userId,
		Email:        email,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return nil, errors.Join(errors.New("could not create user"), err)
	}

	accessToken, err := jwt.NewAccessToken(userId)
	if err != nil {
		return nil, errors.Join(errors.New("could not create access token"), err)
	}
	refreshToken, err := jwt.NewRefreshToken(userId)
	if err != nil {
		return nil, errors.Join(errors.New("could not create refresh token"), err)
	}

	return &resToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) DevSignIn(ctx context.Context, email, password string) (*resToken, error) {
	res, err := s.store.Auth.GetDevUser(ctx, email)
	if err != nil {
		return nil, errors.Join(errors.New("could not get user"), err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.Join(errors.New("wrong password"), err)
	}

	accessToken, err := jwt.NewAccessToken(res.UserId)
	if err != nil {
		return nil, errors.Join(errors.New("could not create access token"), err)
	}
	refreshToken, err := jwt.NewRefreshToken(res.UserId)
	if err != nil {
		return nil, errors.Join(errors.New("could not create refresh token"), err)
	}

	return &resToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
