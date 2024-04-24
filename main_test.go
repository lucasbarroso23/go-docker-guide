package main

import (
	"testing"
)

func Test_intMin(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "positive test",
			args: args{a: 4, b: 5},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intMin(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("intMin() = %v, want %v", got, tt.want)
			}
		})
	}
}
