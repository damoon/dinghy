package middleware

import (
	"testing"
)

func Test_randString(t *testing.T) {
	type args struct {
		r func() int64
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				func() int64 { return 0 },
				0,
			},
			want: "",
		},
		{
			name: "one call",
			args: args{
				func() int64 { return 0 },
				5,
			},
			want: "00000",
		},
		{
			name: "two calls",
			args: args{
				func() int64 { return 0 },
				30,
			},
			want: "000000000000000000000000000000",
		},
		{
			name: "2019 one call",
			args: args{
				func() int64 { return 2019 },
				15,
			},
			want: "0000000000007e3",
		},
		{
			name: "2019 two calls",
			args: args{
				func() int64 { return 2019 },
				30,
			},
			want: "0000000000007e30000000000007e3",
		},
		{
			name: "2019 three calls",
			args: args{
				func() int64 { return 2019 },
				32,
			},
			want: "e30000000000007e30000000000007e3",
		},
		{
			name: "20192019",
			args: args{
				func() int64 { return 20192019 },
				32,
			},
			want: "13000000001341b13000000001341b13",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randString(tt.args.r, tt.args.n); got != tt.want {
				t.Errorf("randString() = %v, want %v", got, tt.want)
			}
		})
	}
}
