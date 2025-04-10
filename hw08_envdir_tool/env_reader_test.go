package main

import (
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    Environment
		wantErr bool
	}{
		{
			name: "case 01",
			args: args{dir: "./testdata/env"},
			want: Environment{
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"FOO": EnvValue{Value: `   foo
with new line`, NeedRemove: false},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
				"EMPTY": EnvValue{Value: "", NeedRemove: true},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}
