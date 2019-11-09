package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type storageMock struct {
	existsValue  bool
	existsError  error
	presignValue *url.URL
	presignError error
}

func (s storageMock) exists(ctx context.Context, objectName string) (bool, error) {
	return s.existsValue, s.existsError
}
func (s storageMock) presign(method string, objectName string) (*url.URL, error) {
	return s.presignValue, s.presignError
}

func TestNewPresignHandler(t *testing.T) {
	type args struct {
		storage     presignStorage
		redirectURL string
		req         *http.Request
	}
	type want struct {
		httpStatus int
		targetURL  string
	}
	type test struct {
		name string
		args args
		want want
	}
	testcase := func(url string) test {
		return test{
			name: "redirect root",
			args: args{
				storage:     storageMock{},
				redirectURL: "http://mockminio:9000/",
				req:         httptest.NewRequest("GET", url, nil),
			},
			want: want{
				httpStatus: http.StatusTemporaryRedirect,
				targetURL:  "http://mockminio:9000/",
			},
		}
	}
	tests := []test{
		testcase("/"),
		testcase("/file"),
		testcase("/some/missing/path"),
		testcase("/directory"),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewPresignHandler(tt.args.storage, tt.args.redirectURL)
			w := httptest.NewRecorder()
			h(w, tt.args.req)
			got := w.Result()
			defer got.Body.Close()
			if !reflect.DeepEqual(got.StatusCode, tt.want.httpStatus) {
				t.Errorf("presign Handler returned http status code %v, want %v", got.StatusCode, tt.want.httpStatus)
			}
			l := got.Header.Get("Location")
			if !reflect.DeepEqual(l, tt.want.targetURL) {
				t.Errorf("presign Handler returned http status code %v, want %v", l, tt.want.targetURL)
			}
		})
	}
}
