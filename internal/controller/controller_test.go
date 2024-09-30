package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/irbgeo/usdt-rate/internal/controller"
	"github.com/irbgeo/usdt-rate/internal/controller/mocks"
	"github.com/irbgeo/usdt-rate/internal/utils/logging"
	rateerr "github.com/irbgeo/usdt-rate/internal/utils/rate-error"
)

func TestController_Rate(t *testing.T) {
	testCases := []struct {
		name              string
		input             controller.Pair
		mockRateGet       func(*mocks.RateProvider)
		mockStorageCreate func(*mocks.Storage)
		mockMetricsCreate func(*mocks.Metrics)
		expectedRate      *controller.Rate
		expectedError     error
	}{
		{
			name:  "success",
			input: controller.Pair{TokenB: "rub"},
			mockRateGet: func(m *mocks.RateProvider) {
				m.On("Get", mock.Anything, controller.Pair{TokenA: "usdt", TokenB: "rub"}).
					Return(&controller.Rate{Ask: "75.5", Bid: "74.5"}, nil)
			},
			mockStorageCreate: func(m *mocks.Storage) {
				m.On("Create", mock.Anything, controller.Rate{Ask: "75.5", Bid: "74.5"}).
					Return(nil)
			},
			mockMetricsCreate: func(m *mocks.Metrics) {
				m.On("ObserveControllerRate", controller.Pair{TokenB: "rub"}, mock.Anything, nil)
			},
			expectedRate: &controller.Rate{Ask: "75.5", Bid: "74.5"},
		},
		{
			name:  "rate_provider_error",
			input: controller.Pair{TokenB: "eur"},
			mockRateGet: func(m *mocks.RateProvider) {
				m.On("Get", mock.Anything, controller.Pair{TokenA: "usdt", TokenB: "eur"}).
					Return(nil, errors.New("rate provider error"))
			},
			mockStorageCreate: func(m *mocks.Storage) {},
			mockMetricsCreate: func(m *mocks.Metrics) {
				m.On("ObserveControllerRate", controller.Pair{TokenB: "eur"}, mock.Anything, rateerr.New(errors.New("rate provider error"), "msg", "failed to get rate"))
			},
			expectedError: rateerr.New(errors.New("rate provider error"), "msg", "failed to get rate"),
		},
		{
			name:  "storage_error",
			input: controller.Pair{TokenB: "gbp"},
			mockRateGet: func(m *mocks.RateProvider) {
				m.On("Get", mock.Anything, controller.Pair{TokenA: "usdt", TokenB: "gbp"}).
					Return(&controller.Rate{Ask: "0.8", Bid: "0.79"}, nil)
			},
			mockStorageCreate: func(m *mocks.Storage) {
				m.On("Create", mock.Anything, controller.Rate{Ask: "0.8", Bid: "0.79"}).
					Return(errors.New("storage error"))
			},
			mockMetricsCreate: func(m *mocks.Metrics) {
				m.On("ObserveControllerRate", controller.Pair{TokenB: "gbp"}, mock.Anything, rateerr.New(errors.New("storage error"), "msg", "failed to store rate"))
			},
			expectedError: rateerr.New(errors.New("storage error"), "msg", "failed to store rate"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logging.Init("test")
			// Setup mocks
			mockRateProvider := mocks.NewRateProvider(t)
			tc.mockRateGet(mockRateProvider)

			mockStorage := mocks.NewStorage(t)
			tc.mockStorageCreate(mockStorage)

			mockMetrics := mocks.NewMetrics(t)
			tc.mockMetricsCreate(mockMetrics)

			// Create controller
			c := controller.NewService(mockRateProvider, mockStorage, mockMetrics)

			// Call the method
			rate, err := c.Rate(context.Background(), tc.input)

			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedRate, rate)
		})
	}
}
