package utils

import (
	"os"
	"reflect"
	"testing"
)

func TestWaitForSignal(t *testing.T) {
	type args struct {
		sources []os.Signal
	}
	tests := []struct {
		name string
		args args
		want os.Signal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WaitForSignal(tt.args.sources...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WaitForSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}
