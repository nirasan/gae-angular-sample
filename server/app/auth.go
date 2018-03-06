package app

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"os"
	"strings"
)

func GetHMACKey() (uuid.UUID, error) {
	key := os.Getenv("HMAC_KEY")
	if key == "" {
		return uuid.NewV4(), nil
	}
	return uuid.FromString(key)
}

func GetTokenFromRequest(r *http.Request) (*jwt.Token, error) {
	h := r.Header.Get("Authorization")

	if h == "" {
		return nil, errors.New("Auth header empty")
	}

	parts := strings.SplitN(h, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("Invalid auth header")
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		id, err := GetHMACKey()
		if err != nil {
			return nil, err
		}
		return id.Bytes(), nil
	})
}

func AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		ctx := appengine.NewContext(e.Request())
		token, err := GetTokenFromRequest(e.Request())
		if err != nil {
			log.Errorf(ctx, "Get token: %v", err)
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			log.Errorf(ctx, "Get Claims: %v", err)
			return echo.NewHTTPError(http.StatusUnauthorized)
		} else if sub, ok := claims["sub"].(string); !ok {
			log.Errorf(ctx, "Get Sub: %v", err)
			return echo.NewHTTPError(http.StatusUnauthorized)
		} else {
			e.Set("UserID", sub)
		}
		return next(e)
	}
}
