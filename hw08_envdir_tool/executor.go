package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	exc := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	envStr := make([]string, len(env))
	i := 0
	for k, v := range env {
		envStr[i] = fmt.Sprintf("%s=%s", k, v.Value)
		i++
	}
	exc.Env = envStr
	slog.Info("Строка переменных для запуска", "envStr", envStr)

	exc.Stdout = os.Stdout
	exc.Stdin = os.Stdin
	err := exc.Run()
	if err != nil {
		slog.Error("exc.Run", "Error", err)
		returnCode = 1
	}
	return
}
