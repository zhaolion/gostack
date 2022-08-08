package stringutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompact(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{a: []string{"a", ""}},
			want: []string{"a"},
		},
		{
			args: args{a: []string{""}},
			want: []string{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Compact(tt.args.a))
		})
	}
}

func TestContain(t *testing.T) {
	type args struct {
		list []string
		s    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				list: []string{"a", "b"},
				s:    "a",
			},
			want: true,
		},
		{
			args: args{
				list: []string{"a", "b"},
				s:    "c",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contain(tt.args.list, tt.args.s); got != tt.want {
				t.Errorf("Contain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type args struct {
		all []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{
				all: []string{"c", "a", "a", "b"},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.args.all); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualIgnoreCase(t *testing.T) {
	tests := []struct {
		name string
		lhs  string
		rhs  string
		want bool
	}{
		{
			name: "empty",
			lhs:  "",
			rhs:  "",
			want: true,
		}, {
			name: "equal",
			lhs:  "hello",
			rhs:  "hello",
			want: true,
		}, {
			name: "equal ignore case",
			lhs:  "hello",
			rhs:  "hELlO",
			want: true,
		}, {
			name: "not equal",
			lhs:  "hello",
			rhs:  "world",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualIgnoreCase(tt.lhs, tt.rhs); got != tt.want {
				t.Errorf("EqualIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
