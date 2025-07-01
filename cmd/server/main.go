package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/configs"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/configs/memory"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/handlers"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/models"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/repositories"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/services"
)

func main() {
	err := command()
	if err != nil {
		panic(err)
	}
}

func command() error {
	config := parseFlags()

	srv, err := newServer(config)
	if err != nil {
		return err
	}

	err = runServer(
		context.Background(),
		srv,
	)
	if err != nil {
		return err
	}

	return nil
}

func parseFlags() *configs.ServerConfig {
	return configs.NewServerConfig()
}

func newServer(
	config *configs.ServerConfig,
) (*http.Server, error) {
	memStorage := memory.NewMemory[models.MetricID, models.Metrics]()

	metricsMemoryGetRepository := repositories.NewMetricsMemoryGetRepository(memStorage)
	metricsMemorySaveRepository := repositories.NewMetricsMemorySaveRepository(memStorage)

	metricUpdateService := services.NewMetricUpdateService(
		services.WithMetricUpdateGetter(metricsMemoryGetRepository),
		services.WithMetricUpdateSaver(metricsMemorySaveRepository),
	)

	metricUpdateHandler := handlers.NewMetricUpdatePathHandler(
		handlers.WithMetricUpdaterPath(metricUpdateService),
	)

	router := chi.NewRouter()

	metricUpdateHandler.RegisterRoute(router)

	srv := &http.Server{Addr: config.Address, Handler: router}

	return srv, nil
}

func runServer(
	ctx context.Context,
	srv *http.Server,
) error {
	ctx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	errChan := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}
	case err := <-errChan:
		if err != nil {
			return err
		}
	}

	return nil
}
