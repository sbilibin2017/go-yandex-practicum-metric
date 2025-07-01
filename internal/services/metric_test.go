package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestMetricUpdateService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := NewMockGetter(ctrl)
	mockSaver := NewMockSaver(ctrl)

	svc := NewMetricUpdateService(
		WithMetricUpdateGetter(mockGetter),
		WithMetricUpdateSaver(mockSaver),
	)

	ctx := context.Background()

	// Helper to create pointer to int64
	int64Ptr := func(v int64) *int64 { return &v }
	float64Ptr := func(v float64) *float64 { return &v }

	tests := []struct {
		name           string
		inputMetrics   []*models.Metrics
		mockGetterFunc func()
		mockSaverFunc  func() error
		expected       []*models.Metrics
		expectErr      bool
	}{
		{
			name: "Update counter metric with existing delta",
			inputMetrics: []*models.Metrics{
				{
					ID:    "metric1",
					MType: models.Counter,
					Delta: int64Ptr(5),
				},
			},
			mockGetterFunc: func() {
				mockGetter.EXPECT().
					Get(ctx, models.MetricID{ID: "metric1", MType: models.Counter}).
					Return(&models.Metrics{ID: "metric1", MType: models.Counter, Delta: int64Ptr(10)}, nil)
			},
			mockSaverFunc: func() error {
				mockSaver.EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil)
				return nil
			},
			expected: []*models.Metrics{
				{
					ID:    "metric1",
					MType: models.Counter,
					Delta: int64Ptr(15), // 5 + 10
				},
			},
		},
		{
			name: "Save gauge metric without getter call",
			inputMetrics: []*models.Metrics{
				{
					ID:    "metric2",
					MType: models.Gauge,
					Value: float64Ptr(3.14),
				},
			},
			mockGetterFunc: func() {},
			mockSaverFunc: func() error {
				mockSaver.EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil)
				return nil
			},
			expected: []*models.Metrics{
				{
					ID:    "metric2",
					MType: models.Gauge,
					Value: float64Ptr(3.14),
				},
			},
		},
		{
			name: "Getter returns error",
			inputMetrics: []*models.Metrics{
				{
					ID:    "metric3",
					MType: models.Counter,
					Delta: int64Ptr(1),
				},
			},
			mockGetterFunc: func() {
				mockGetter.EXPECT().
					Get(ctx, models.MetricID{ID: "metric3", MType: models.Counter}).
					Return(nil, errors.New("getter error"))
			},
			mockSaverFunc: func() error {
				return nil // won't be called
			},
			expectErr: true,
		},
		{
			name: "Saver returns error",
			inputMetrics: []*models.Metrics{
				{
					ID:    "metric4",
					MType: models.Gauge,
					Value: float64Ptr(2.71),
				},
			},
			mockGetterFunc: func() {},
			mockSaverFunc: func() error {
				mockSaver.EXPECT().
					Save(ctx, gomock.Any()).
					Return(errors.New("save error"))
				return nil
			},
			expectErr: true,
		},
		{
			name: "Skip nil metrics in input",
			inputMetrics: []*models.Metrics{
				nil,
				{
					ID:    "metric5",
					MType: models.Gauge,
					Value: float64Ptr(1.23),
				},
			},
			mockGetterFunc: func() {},
			mockSaverFunc: func() error {
				mockSaver.EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil)
				return nil
			},
			expected: []*models.Metrics{
				{
					ID:    "metric5",
					MType: models.Gauge,
					Value: float64Ptr(1.23),
				},
			},
		},
		{
			name: "metrics returned sorted by ID",
			inputMetrics: []*models.Metrics{
				{
					ID:    "z_metric",
					MType: models.Gauge,
					Value: float64Ptr(1.1),
				},
				{
					ID:    "a_metric",
					MType: models.Gauge,
					Value: float64Ptr(2.2),
				},
				{
					ID:    "m_metric",
					MType: models.Gauge,
					Value: float64Ptr(3.3),
				},
			},
			mockGetterFunc: func() {},
			mockSaverFunc: func() error {
				mockSaver.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(3)
				return nil
			},
			expected: []*models.Metrics{
				{
					ID:    "a_metric",
					MType: models.Gauge,
					Value: float64Ptr(2.2),
				},
				{
					ID:    "m_metric",
					MType: models.Gauge,
					Value: float64Ptr(3.3),
				},
				{
					ID:    "z_metric",
					MType: models.Gauge,
					Value: float64Ptr(1.1),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockGetterFunc()
			tt.mockSaverFunc()

			got, err := svc.Update(ctx, tt.inputMetrics)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(got))

			for i := range got {
				assert.Equal(t, tt.expected[i].ID, got[i].ID)
				assert.Equal(t, tt.expected[i].MType, got[i].MType)

				if tt.expected[i].Delta == nil {
					assert.Nil(t, got[i].Delta)
				} else {
					assert.NotNil(t, got[i].Delta)
					assert.Equal(t, *tt.expected[i].Delta, *got[i].Delta)
				}

				if tt.expected[i].Value == nil {
					assert.Nil(t, got[i].Value)
				} else {
					assert.NotNil(t, got[i].Value)
					assert.Equal(t, *tt.expected[i].Value, *got[i].Value)
				}
			}
		})
	}
}
