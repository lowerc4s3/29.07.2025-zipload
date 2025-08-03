package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	HttpErrMalformedRequest = echo.NewHTTPError(http.StatusBadRequest, "malformed body")
	HttpErrMalformedID      = echo.NewHTTPError(http.StatusBadRequest, "malformed ID")
	HttpErrDownload         = echo.NewHTTPError(http.StatusBadGateway, "all requested resources returned error")
	HttpErrForbiddenMime    = echo.NewHTTPError(http.StatusBadRequest, "all files had forbidden type")
	HttpErrTaskNotFound     = echo.NewHTTPError(http.StatusNotFound, "no task found with provided ID")
	HttpErrTooManyTasks     = echo.NewHTTPError(http.StatusServiceUnavailable, "server is busy")
	HttpErrTooManyTaskFiles = echo.NewHTTPError(http.StatusBadRequest, "too many files requested for the task")
	HttpErrTooManySources   = echo.NewHTTPError(http.StatusBadRequest, "too many sources requested")
)
