package repositories

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/configs/memory"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/models"
)

func TestMetricsMemorySaveRepository_Save(t *testing.T) {
	mem := memory.NewMemory[models.MetricID, models.Metrics]()
	repo := NewMetricsMemorySaveRepository(mem)
	ctx := context.Background()

	metric1 := models.Metrics{ID: "metric1", MType: "gauge"}
	metric2 := models.Metrics{ID: "metric2", MType: "counter"}

	err := repo.Save(ctx, metric1)
	require.NoError(t, err)

	err = repo.Save(ctx, metric2)
	require.NoError(t, err)

	mem.Mu.RLock()
	defer mem.Mu.RUnlock()

	assert.Len(t, mem.Data, 2)
	assert.Equal(t, metric1, mem.Data[models.MetricID{ID: "metric1", MType: "gauge"}])
	assert.Equal(t, metric2, mem.Data[models.MetricID{ID: "metric2", MType: "counter"}])
}

func TestMetricsMemorySaveRepository_Save_Concurrent(t *testing.T) {
	mem := memory.NewMemory[models.MetricID, models.Metrics]()
	repo := NewMetricsMemorySaveRepository(mem)
	ctx := context.Background()

	metric1 := models.Metrics{ID: "metric1", MType: "gauge"}
	metric2 := models.Metrics{ID: "metric2", MType: "counter"}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := repo.Save(ctx, metric1)
		require.NoError(t, err)
	}()

	go func() {
		defer wg.Done()
		err := repo.Save(ctx, metric2)
		require.NoError(t, err)
	}()

	wg.Wait()

	mem.Mu.RLock()
	defer mem.Mu.RUnlock()

	assert.Len(t, mem.Data, 2)
	assert.Contains(t, mem.Data, models.MetricID{ID: "metric1", MType: "gauge"})
	assert.Contains(t, mem.Data, models.MetricID{ID: "metric2", MType: "counter"})
}

func TestMetricsMemoryGetRepository_Get_Found(t *testing.T) {
	mem := memory.NewMemory[models.MetricID, models.Metrics]()

	metric := models.Metrics{ID: "metric1", MType: "gauge"}
	key := models.MetricID{ID: metric.ID, MType: metric.MType}

	mem.Data[key] = metric

	repo := NewMetricsMemoryGetRepository(mem)
	ctx := context.Background()

	got, err := repo.Get(ctx, key)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, metric, *got)
}

func TestMetricsMemoryGetRepository_Get_NotFound(t *testing.T) {
	mem := memory.NewMemory[models.MetricID, models.Metrics]()
	repo := NewMetricsMemoryGetRepository(mem)
	ctx := context.Background()

	key := models.MetricID{ID: "missing", MType: "counter"}

	got, err := repo.Get(ctx, key)

	require.NoError(t, err)
	assert.Nil(t, got)
}
