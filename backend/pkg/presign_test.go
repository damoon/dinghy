package dinghy

import (
	"context"
	"fmt"
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

	redirectTestcase := func(n, u string) test {
		return test{
			name: n,
			args: args{
				storage:     storageMock{},
				redirectURL: "http://mockminio:9000/",
				req:         httptest.NewRequest("GET", u, nil),
			},
			want: want{
				httpStatus: http.StatusTemporaryRedirect,
				targetURL:  "http://mockminio:9000/",
			},
		}
	}
	downloadTestcase := func(n, u string) test {
		ur, err := url.Parse(fmt.Sprintf("%s%s", "http://mockminio:9000/", u))
		if err != nil {
			panic(err)
		}

		return test{
			name: n,
			args: args{
				storage: storageMock{
					existsValue:  true,
					presignValue: ur,
				},
				req: httptest.NewRequest("GET", u, nil),
			},
			want: want{
				httpStatus: http.StatusTemporaryRedirect,
				targetURL:  ur.String(),
			},
		}
	}
	uploadTestcase := func(n, u string) test {
		ur, err := url.Parse(fmt.Sprintf("%s%s", "http://mockminio:9000/", u))
		if err != nil {
			panic(err)
		}

		return test{
			name: n,
			args: args{
				storage: storageMock{presignValue: ur},
				req:     httptest.NewRequest("PUT", u, nil),
			},
			want: want{
				httpStatus: http.StatusTemporaryRedirect,
				targetURL:  ur.String(),
			},
		}
	}
	tests := []test{
		redirectTestcase("download missing root", "/"),
		redirectTestcase("download missing file", "/file"),
		redirectTestcase("download missing path", "/some/missing/path"),
		redirectTestcase("download missing directory", "/directory/"),
		downloadTestcase("download existing root", "/"),
		downloadTestcase("download existing file", "/file"),
		downloadTestcase("download existing path", "/some/missing/path"),
		downloadTestcase("download existing directory", "/directory/"),
		uploadTestcase("upload root", "/"),
		uploadTestcase("upload file", "/file"),
		uploadTestcase("upload path", "/some/missing/path"),
		uploadTestcase("upload directory", "/directory/"),
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
