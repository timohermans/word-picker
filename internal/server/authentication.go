package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"context"

	"github.com/labstack/echo/v5"

	"github.com/markbates/goth/gothic"
)

type User struct {
	Name      string
	Id        string
	ExpiresAt time.Time
}

const (
	UrlAuthenticate  = "/auth"
	UrlAuthCallback  = "/auth/callback"
	UrlLogout        = "/auth/logout"
	gothProviderName = "openid-connect"
)

func NewAuthenticationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			request := c.Request()

			if strings.Contains(request.RequestURI, "auth") {
				return next(c)
			}

			user, err := GetUserFromSession(c)

			if err != nil {
				return fmt.Errorf("Getting the user from session: %w", err)
			}

			now := time.Now().UTC()
			isExpired := user != nil && now.After(user.ExpiresAt)
			if user == nil || isExpired {
				c.Logger().Info("Anonymous user or expired. Redirecting to auth.", "isAnonymous", user == nil, "isExpired", isExpired)
				redirectStatus := http.StatusFound
				if request.Method != "GET" {
					redirectStatus = http.StatusTemporaryRedirect
				}

				return c.Redirect(redirectStatus, fmt.Sprintf("%s?next=%s", UrlAuthenticate, request.RequestURI))
			}

			return next(c)
		}
	}
}

func GetUserFromSession(c *echo.Context) (*User, error) {
	session, err := GetSession(c, "auth")
	if err != nil {
		return nil, fmt.Errorf("Getting the auth session: %w", err)
	}

	// TODO: unit test this with nil values
	username, ok1 := session.Values["name"].(string)
	userId, ok2 := session.Values["userId"].(string)
	expiresAt, ok3 := session.Values["expiresAt"].(string)
	if !ok1 || !ok2 || !ok3 {
		return nil, nil
	}

	c.Logger().Info("Getting the username from the session", "user", username)

	expiresAtTime, err := time.Parse("2006-01-02 15:04:05", expiresAt)
	if err != nil {
		return nil, fmt.Errorf("Parsing the expiresat %s time: %w", expiresAt, err)
	}

	return &User{
		Name:      username,
		Id:        userId,
		ExpiresAt: expiresAtTime,
	}, nil
}

func RegisterAuthentication(e *echo.Echo, defaultNextUrl string) {
	e.GET(UrlAuthenticate, func(c *echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, gothProviderName)
		request := c.Request().WithContext(ctx)
		response := c.Response()

		session, err := GetSession(c, "auth")
		if err != nil {
			return fmt.Errorf("Getting the auth session store: %w", err)
		}

		next := c.QueryParamOr("next", defaultNextUrl)
		session.Values["next"] = next
		if err := session.Save(request, response); err != nil {
			return fmt.Errorf("Saving previous url in session: %w", err)
		}

		if gothUser, err := gothic.CompleteUserAuth(response, request); err == nil {
			c.Logger().Info("User logged in", "user", gothUser.Name)
			return c.Redirect(http.StatusSeeOther, next)
		}

		gothic.BeginAuthHandler(response, request)

		return nil
	})

	e.GET(UrlAuthCallback, func(c *echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), gothic.ProviderParamKey, gothProviderName)
		request := c.Request().WithContext(ctx)
		response := c.Response()

		session, err := GetSession(c, "auth")
		if err != nil {
			return fmt.Errorf("Getting the auth session store: %w", err)
		}

		user, err := gothic.CompleteUserAuth(response, request)
		if err != nil {
			return fmt.Errorf("Handling the auth callback: %w", err)
		}

		next, ok := session.Values["next"].(string)
		if !ok {
			next = defaultNextUrl
		}
		session.Values["next"] = nil

		session.Values["name"] = user.Name
		session.Values["userId"] = user.UserID
		session.Values["expiresAt"] = user.ExpiresAt.Format("2006-01-02 15:04:05")
		if err := session.Save(request, response); err != nil {
			return fmt.Errorf("Saving user info in session: %w", err)
		}

		c.Logger().Info("User logged in", "user", user.Name)
		return c.Redirect(http.StatusSeeOther, next)
	})
}
