package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/lowerc4s3/29.07.2025-zipload/app"
	"github.com/lowerc4s3/29.07.2025-zipload/batch"
	"github.com/lowerc4s3/29.07.2025-zipload/task"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	batchService := batch.NewBatchService(
		// TODO: Impl config file
		[]string{
			"image/jpeg",
			"application/pdf",
		},
		3,
	)
	taskService := new(task.TaskService)
	handler := app.NewHandler(batchService, taskService)

	app := app.NewApp(handler)
	runApp(app, "8080")
}

func runApp(app *echo.Echo, port string) {
	appCtx, appStop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer appStop()

	go func() {
		if err := app.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server couldn't be recovered")
		}
	}()
	<-appCtx.Done()

	log.Info().Msg("Received SIGINT, shutting down...")
	log.Info().Msg("Press Ctrl-C to shutdown immediately")

	stopCtx, stopCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopCancel()
	if err := app.Shutdown(stopCtx); err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown gracefully")
	}
}
