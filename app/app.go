package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewApp(h *Handler) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(newZerologMiddleware(&log.Logger))
	setupRoutes(app, h)
	return app
}

func newZerologMiddleware(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info().
					Str("method", v.Method).
					Str("URI", v.URI).
					Int("status", v.Status).
					Msg("Request success")
			} else {
				logger.Error().
					Err(v.Error).
					Str("method", v.Method).
					Str("URI", v.URI).
					Int("status", v.Status).
					Msg("Request error")
			}
			return nil
		},
	})
}

func setupRoutes(app *echo.Echo, h *Handler) {
	app.POST("/downloadBatch", h.DownloadBatch)

	app.POST("/createTask", h.CreateTask)
	app.POST("/appendTask", h.AppendTask)
	app.GET("/checkTask/:taskid", h.CheckTask)
	app.GET("/downloadTask/:taskid", h.DownloadTask).Name = "downloadTask"
}
