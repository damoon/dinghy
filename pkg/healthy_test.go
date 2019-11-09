package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type healthMock struct {
	e error
}

func (h healthMock) healthy(ctx context.Context) error {
	return h.e
}

func TestHealthHandler(t *testing.T) {
	type args struct {
		h   healthy
		req *http.Request
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "healthy",
			args: args{
				h:   healthMock{e: nil},
				req: httptest.NewRequest("GET", "/health", nil),
			},
			want: http.StatusOK,
		},
		{
			name: "unhealthy",
			args: args{
				h:   healthMock{e: fmt.Errorf("error to fail test")},
				req: httptest.NewRequest("GET", "/health", nil),
			},
			want: http.StatusServiceUnavailable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := HealthHandler(tt.args.h)
			w := httptest.NewRecorder()
			h(w, tt.args.req)
			resp := w.Result()
			defer resp.Body.Close()
			got := resp.StatusCode
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("healthHandler returned http status code %v, want %v", got, tt.want)
			}
		})
	}
}
