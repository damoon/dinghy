package dinghy

import "testing"

func Test_shouldRedirect(t *testing.T) {
	tests := []struct {
		name     string
		rawQuery string
		want     bool
		wantErr  bool
	}{
		{
			name:     "default",
			rawQuery: "",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "redirect",
			rawQuery: "redirect",
			want:     true,
		},
		{
			name:     "redirect=",
			rawQuery: "redirect=",
			want:     true,
		},
		{
			name:     "redirect=1",
			rawQuery: "redirect=1",
			want:     true,
		},
		{
			name:     "redirect=0",
			rawQuery: "redirect=0",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shouldRedirect(tt.rawQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("shouldRedirect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("shouldRedirect() = %v, want %v", got, tt.want)
			}
		})
	}
}
