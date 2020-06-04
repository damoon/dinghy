package dinghy

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
		err error
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
				err: nil,
				req: httptest.NewRequest("GET", "/healthz", nil),
			},
			want: http.StatusOK,
		},
		{
			name: "unhealthy",
			args: args{
				err: fmt.Errorf("error to fail test"),
				req: httptest.NewRequest("GET", "/healthz", nil),
			},
			want: http.StatusServiceUnavailable,
		},
	}
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAdminServer()
			s.Storage = &healthMock{e: tt.args.err}

			w := httptest.NewRecorder()
			s.handleHealthz()(w, tt.args.req)

			resp := w.Result()
			defer resp.Body.Close()

			got := resp.StatusCode
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("healthHandler returned http status code %v, want %v", got, tt.want)
			}
		})
	}
}
