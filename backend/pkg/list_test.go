package dinghy

import (
	"testing"
)

func Test_parseRequest(t *testing.T) {
	tests := []struct {
		name     string
		rawQuery string
		want     bool
		want1    bool
		wantErr  bool
	}{
		{
			name: "redirect not set",
		},
		{
			name:     "redirect set without value",
			rawQuery: "redirect",
			want:     true,
		},
		{
			name:     "redirect set to empty string",
			rawQuery: "redirect=",
			want:     true,
		},
		{
			name:     "redirect set to 1",
			rawQuery: "redirect=1",
			want:     true,
		},
		{
			name:     "redirect set to 0",
			rawQuery: "redirect=0",
			want:     true,
		},
		{
			name:     "thumbnail set without value",
			rawQuery: "thumbnail",
			want1:    true,
		},
		{
			name:     "thumbnail set to empty string",
			rawQuery: "thumbnail=",
			want1:    true,
		},
		{
			name:     "thumbnail set to 1",
			rawQuery: "thumbnail=1",
			want1:    true,
		},
		{
			name:     "thumbnail set to 0",
			rawQuery: "thumbnail=0",
			want1:    true,
		},
		{
			name:     "redirect and thumbnail set",
			rawQuery: "redirect&thumbnail",
			want:     true,
			want1:    true,
		},
		{
			name:     "thumbnail and redirect set",
			rawQuery: "thumbnail&redirect",
			want:     true,
			want1:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseRequest(tt.rawQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseRequest() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
