package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrMalformedRequest = echo.NewHTTPError(http.StatusBadRequest, "malformed body")
	ErrMalformedID      = echo.NewHTTPError(http.StatusBadRequest, "malformed ID")
	ErrDownload         = echo.NewHTTPError(http.StatusBadGateway, "all requested resources returned error")
	ErrForbiddenMime    = echo.NewHTTPError(http.StatusBadRequest, "all files had forbidden type")
	ErrTaskNotFound     = echo.NewHTTPError(http.StatusNotFound, "no task found with provided ID")
	ErrTooManyTasks     = echo.NewHTTPError(http.StatusServiceUnavailable, "tasks' max amount achieved")
)
