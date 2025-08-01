package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewApp() *echo.Echo {
	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(newZerologMiddleware(&log.Logger))
	setupRoutes(app)
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

func setupRoutes(app *echo.Echo) {
	app.POST("/downloadBatch", func(c echo.Context) error { panic("todo") })

	app.POST("/createTask", func(c echo.Context) error { panic("todo") })
	app.POST("/appendTask", func(c echo.Context) error { panic("todo") })
	app.GET("/checkTask/:taskid", func(c echo.Context) error { panic("todo") })
	app.GET("/downloadTask/:taskid", func(c echo.Context) error { panic("todo") })
}
