package main

import (
	"log/slog"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	err := checkVariables(env)
	if err != nil {
		slog.Error("Ошибка обработки переменных окружения", "checkVariables", err)
		returnCode = 1
	}
	exc := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	envStr := os.Environ()
	exc.Env = envStr
	exc.Stdout = os.Stdout
	exc.Stdin = os.Stdin
	err = exc.Run()
	if err != nil {
		slog.Error("exc.Run", "Error", err)
		returnCode = 1
	}
	return
}
