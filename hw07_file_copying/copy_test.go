package main

import (
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	type args struct {
		fromPath string
		toPath   string
		offset   int64
		limit    int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case 01: offset exceeds file size",
			args: args{
				fromPath: "testdata/input.txt",
				toPath:   "testdata/out.txt",
				offset:   10000,
				limit:    0,
			},
			wantErr: true,
		},
		{
			name: "case 02: copy ok",
			args: args{
				fromPath: "testdata/input.txt",
				toPath:   "testdata/out.txt",
				offset:   0,
				limit:    0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit); (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
			os.Remove(tt.args.toPath)
		})
	}
}
