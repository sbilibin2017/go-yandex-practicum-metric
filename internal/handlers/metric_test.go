package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-yandex-practicum-metric/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := NewMockMetricUpdater(ctrl)
	handler := NewMetricUpdatePathHandler(WithMetricUpdaterPath(mockUpdater))

	r := chi.NewRouter()
	handler.RegisterRoute(r)

	tests := []struct {
		name         string
		method       string
		url          string
		mockExpect   func()
		expectedCode int
	}{
		{
			name:   "Valid counter metric",
			method: http.MethodPost,
			url:    "/update/counter/myCounter/100",
			mockExpect: func() {
				mockUpdater.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf([]*models.Metrics{})).
					Return(nil, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "Valid gauge metric",
			method: http.MethodPost,
			url:    "/update/gauge/myGauge/99.9",
			mockExpect: func() {
				mockUpdater.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf([]*models.Metrics{})).
					Return(nil, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Missing metric name",
			method:       http.MethodPost,
			url:          "/update/counter//100", // empty name param triggers 404
			mockExpect:   func() {},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Missing metric value",
			method:       http.MethodPost,
			url:          "/update/gauge/myGauge", // no value param triggers 400
			mockExpect:   func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid counter value",
			method:       http.MethodPost,
			url:          "/update/counter/myCounter/notanumber",
			mockExpect:   func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid gauge value",
			method:       http.MethodPost,
			url:          "/update/gauge/myGauge/invalidfloat",
			mockExpect:   func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Unsupported metric type",
			method:       http.MethodPost,
			url:          "/update/unknown/myMetric/123",
			mockExpect:   func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "Updater returns error",
			method: http.MethodPost,
			url:    "/update/counter/myMetric/123",
			mockExpect: func() {
				mockUpdater.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf([]*models.Metrics{})).
					Return(nil, context.DeadlineExceeded)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()

			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
