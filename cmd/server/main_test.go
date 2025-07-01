package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestRunServer_ShutdownOnContextCancel(t *testing.T) {
	srv := &http.Server{
		Addr: "127.0.0.1:0", // use random free port
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error)
	go func() {
		done <- runServer(ctx, srv)
	}()

	// wait for server to start listening
	time.Sleep(100 * time.Millisecond)

	cancel() // cancel context to trigger shutdown

	err := <-done
	require.NoError(t, err)
}

func TestRunServer_ListenAndServeReturnsError(t *testing.T) {
	// create server with invalid address to force ListenAndServe error
	srv := &http.Server{
		Addr: "invalid_addr",
	}

	ctx := context.Background()

	err := runServer(ctx, srv)
	require.Error(t, err)
}

func TestRunServer_ListenAndServeReturnsErrServerClosed(t *testing.T) {
	srv := &http.Server{
		Addr: "127.0.0.1:0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run server in goroutine
	done := make(chan error)
	go func() {
		err := runServer(ctx, srv)
		done <- err
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Gracefully shutdown server directly to simulate http.ErrServerClosed
	err := srv.Shutdown(context.Background())
	require.NoError(t, err)

	// Wait for runServer to return
	err = <-done
	require.NoError(t, err)
}

type ServerSuite struct {
	suite.Suite
	client *resty.Client
}

func (s *ServerSuite) SetupSuite() {
	errChan := make(chan error, 1)

	go func() {
		err := command()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	// Wait a short time for server to start, or receive error
	select {
	case err := <-errChan:
		if err != nil {
			s.T().Fatalf("server failed to start: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		// assume server started OK
	}

	s.client = resty.New().
		SetBaseURL("http://localhost:8080").
		SetTimeout(3 * time.Second)
}

func (s *ServerSuite) TearDownSuite() {
	// Optionally implement graceful shutdown here
	// but without cancel context, might rely on process exit
}

func (s *ServerSuite) TestUpdateMetricScenarios() {
	tests := []struct {
		name       string
		path       string
		pathParams map[string]string
		wantStatus int
	}{
		{
			name:       "success with value",
			path:       "/update/{type}/{name}/{value}",
			pathParams: map[string]string{"type": "counter", "name": "requests", "value": "42"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing metric name",
			path:       "/update/gauge",
			pathParams: map[string]string{},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid metric type",
			path:       "/update/{type}/{name}",
			pathParams: map[string]string{"type": "invalid_type", "name": "temp"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid value",
			path:       "/update/{type}/{name}/{value}",
			pathParams: map[string]string{"type": "counter", "name": "requests", "value": "not_a_number"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := s.client.R()
			if len(tt.pathParams) > 0 {
				req.SetPathParams(tt.pathParams)
			}
			resp, err := req.Post(tt.path)
			s.Require().NoError(err)
			s.Equal(tt.wantStatus, resp.StatusCode())
		})
	}
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
