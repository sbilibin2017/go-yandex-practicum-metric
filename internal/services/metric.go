package services

import (
	"context"
	"sort"

	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/models"
)

type Getter interface {
	Get(ctx context.Context, metricID models.MetricID) (*models.Metrics, error)
}

type Saver interface {
	Save(ctx context.Context, metric models.Metrics) error
}

type MetricUpdateService struct {
	getter Getter
	saver  Saver
}

func NewMetricUpdateService(opts ...MetricUpdateOpt) *MetricUpdateService {
	svc := &MetricUpdateService{}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

type MetricUpdateOpt func(*MetricUpdateService)

func WithMetricUpdateGetter(getter Getter) MetricUpdateOpt {
	return func(svc *MetricUpdateService) {
		svc.getter = getter
	}
}

func WithMetricUpdateSaver(saver Saver) MetricUpdateOpt {
	return func(svc *MetricUpdateService) {
		svc.saver = saver
	}
}

func (svc *MetricUpdateService) Update(
	ctx context.Context,
	metrics []*models.Metrics,
) ([]*models.Metrics, error) {
	updated := make(map[models.MetricID]models.Metrics)

	for _, metric := range metrics {
		if metric == nil {
			continue
		}

		switch metric.MType {
		case models.Counter:
			current, err := svc.getter.Get(ctx, models.MetricID{ID: metric.ID, MType: metric.MType})
			if err != nil {
				return nil, err
			}
			if current != nil && current.Delta != nil && metric.Delta != nil {
				*metric.Delta += *current.Delta
			}
		}

		err := svc.saver.Save(ctx, *metric)
		if err != nil {
			return nil, err
		}

		updated[models.MetricID{ID: metric.ID, MType: metric.MType}] = *metric
	}

	updatedSlice := make([]*models.Metrics, 0, len(updated))
	for _, m := range updated {

		mCopy := m
		updatedSlice = append(updatedSlice, &mCopy)
	}

	sort.Slice(updatedSlice, func(i, j int) bool {
		return updatedSlice[i].ID < updatedSlice[j].ID
	})

	return updatedSlice, nil
}
