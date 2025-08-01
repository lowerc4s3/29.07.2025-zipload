package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrMalformedRequest = echo.NewHTTPError(http.StatusBadRequest, "malformed body")
	ErrMalformedID      = echo.NewHTTPError(http.StatusBadRequest, "malformed ID")
	ErrTaskNotFound     = echo.NewHTTPError(http.StatusNotFound, "no task found with provided ID")
	ErrTooManyTasks     = echo.NewHTTPError(http.StatusServiceUnavailable, "tasks' max amount achieved")
)
