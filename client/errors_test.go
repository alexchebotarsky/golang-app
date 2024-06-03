package client

import (
	"errors"
	"testing"
)

func TestErrMultiple_Error(t *testing.T) {
	type fields struct {
		Errs []error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "no errors",
			fields: fields{Errs: []error{}},
			want:   "",
		},
		{
			name: "one error",
			fields: fields{Errs: []error{
				errors.New("error 1"),
			}},
			want: "error 1",
		},
		{
			name: "two errors",
			fields: fields{Errs: []error{
				errors.New("error 1"),
				errors.New("error 2"),
			}},
			want: "error 1 | error 2",
		},
		{
			name: "three errors",
			fields: fields{Errs: []error{
				errors.New("error 1"),
				errors.New("error 2"),
				errors.New("error 3"),
			}},
			want: "error 1 | error 2 | error 3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrMultiple{
				Errs: tt.fields.Errs,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ErrMultiple.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
