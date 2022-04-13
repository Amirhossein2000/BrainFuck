package brainfuck

import (
	"io"
	"strings"
	"testing"
)

func TestScanInput(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	reader := strings.NewReader("><+-,.[]*")
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "increment pointer",
			args: args{reader: reader},
			want: byte('>'),
		},
		{
			name: "decrement pointer",
			args: args{reader: reader},
			want: byte('<'),
		},
		{
			name: "increment value at pointer",
			args: args{reader: reader},
			want: byte('+'),
		},
		{
			name: "decrement value at pointer",
			args: args{reader: reader},
			want: byte('-'),
		},
		{
			name: "read one character",
			args: args{reader: reader},
			want: byte(','),
		},
		{
			name: "print value at pointer",
			args: args{reader: reader},
			want: byte('.'),
		},
		{
			name: "begin loop",
			args: args{reader: reader},
			want: byte('['),
		},
		{
			name: "end loop",
			args: args{reader: reader},
			want: byte(']'),
		},
		{
			name: "custom operation",
			args: args{reader: reader},
			want: byte('*'),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := scanInput(tt.args.reader); got != tt.want {
				t.Errorf("expected %v(%s) but recieved %v(%s)", got, string(got), tt.want, string(tt.want))
			}
		})
	}
}
