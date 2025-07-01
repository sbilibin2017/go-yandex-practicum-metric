package repositories

import (
	"context"

	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/configs/memory"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/models"
)

type MetricsMemorySaveRepository struct {
	storage *memory.Memory[models.MetricID, models.Metrics]
}

func NewMetricsMemorySaveRepository(
	storage *memory.Memory[models.MetricID, models.Metrics],
) *MetricsMemorySaveRepository {
	return &MetricsMemorySaveRepository{storage: storage}
}

func (r *MetricsMemorySaveRepository) Save(
	ctx context.Context,
	metric models.Metrics,
) error {
	r.storage.Mu.Lock()
	defer r.storage.Mu.Unlock()

	r.storage.Data[models.MetricID{ID: metric.ID, MType: metric.MType}] = metric

	return nil
}

type MetricsMemoryGetRepository struct {
	storage *memory.Memory[models.MetricID, models.Metrics]
}

func NewMetricsMemoryGetRepository(
	storage *memory.Memory[models.MetricID, models.Metrics],
) *MetricsMemoryGetRepository {
	return &MetricsMemoryGetRepository{storage: storage}
}

func (r *MetricsMemoryGetRepository) Get(
	ctx context.Context,
	metricID models.MetricID,
) (*models.Metrics, error) {
	r.storage.Mu.RLock()
	defer r.storage.Mu.RUnlock()

	metric, found := r.storage.Data[models.MetricID{ID: metricID.ID, MType: metricID.MType}]
	if !found {
		return nil, nil
	}

	return &metric, nil
}
