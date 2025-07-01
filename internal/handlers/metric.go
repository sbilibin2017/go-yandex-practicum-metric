package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/models"
)

// MetricUpdater defines an interface for updating multiple metrics.
type MetricUpdater interface {
	Update(ctx context.Context, metrics []*models.Metrics) ([]*models.Metrics, error)
}

// Functional options for MetricUpdatePathHandler
type MetricUpdatePathHandlerOption func(*MetricUpdatePathHandler)

func WithMetricUpdaterPath(svc MetricUpdater) MetricUpdatePathHandlerOption {
	return func(h *MetricUpdatePathHandler) {
		h.svc = svc
	}
}

// MetricUpdatePathHandler handles metric updates via URL path parameters.
type MetricUpdatePathHandler struct {
	svc MetricUpdater
}

func NewMetricUpdatePathHandler(opts ...MetricUpdatePathHandlerOption) *MetricUpdatePathHandler {
	h := &MetricUpdatePathHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *MetricUpdatePathHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var metric models.Metrics
	metric.ID = name

	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metricType {
	case models.Counter:
		delta, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Delta = &delta
		metric.MType = models.Counter

	case models.Gauge:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Value = &val
		metric.MType = models.Gauge

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := h.svc.Update(r.Context(), []*models.Metrics{&metric}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricUpdatePathHandler) RegisterRoute(r chi.Router) {
	r.Post("/update/{type}/{name}/{value}", h.Update)
	r.Post("/update/{type}/{name}", h.Update)
}
