package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"aiqadam-backend/internal/config"
	"aiqadam-backend/internal/http/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server wraps the Echo HTTP server.
type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

// NewServer creates and configures the Echo instance with middleware and routes.
func NewServer(cfg *config.Config, deps routes.Deps) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit("50M"))
	e.Static("/files", cfg.StoragePath)

	corsConfig := middleware.CORSConfig{
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}
	if isDevEnv(cfg.AppEnv) {
		// Любой порт localhost / 127.0.0.1 (Vite 3000, 3001, 5173 и т.д.)
		corsConfig.AllowOriginFunc = func(origin string) (bool, error) {
			return isLocalDevOrigin(origin), nil
		}
	} else if len(cfg.CORSAllowedOrigins) > 0 {
		corsConfig.AllowOrigins = cfg.CORSAllowedOrigins
	} else {
		corsConfig.AllowOrigins = []string{"*"}
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	routes.Register(e, deps)

	return &Server{echo: e, cfg: cfg}
}

// Start listens and serves HTTP requests.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	return s.echo.Start(addr)
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func isDevEnv(appEnv string) bool {
	switch strings.ToLower(strings.TrimSpace(appEnv)) {
	case "development", "dev", "local":
		return true
	default:
		return false
	}
}

func isLocalDevOrigin(origin string) bool {
	if origin == "" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil || u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	return host == "localhost" || host == "127.0.0.1"
}
