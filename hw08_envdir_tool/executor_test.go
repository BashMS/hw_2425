package main

import "testing"

func TestRunCmd(t *testing.T) {
	type args struct {
		cmd []string
		env Environment
	}
	tests := []struct {
		name           string
		args           args
		wantReturnCode int
	}{
		{
			name: "case 01",
			args: args{
				cmd: []string{
					"bin/bash",
					"./testdata/echo.sh",
					"arg1=1",
					"arg2=2",
				},
				env: Environment{
					"one": EnvValue{Value: "1", NeedRemove: false},
				},
			},
			wantReturnCode: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReturnCode := RunCmd(tt.args.cmd, tt.args.env); gotReturnCode != tt.wantReturnCode {
				t.Errorf("RunCmd() = %v, want %v", gotReturnCode, tt.wantReturnCode)
			}
		})
	}
}
