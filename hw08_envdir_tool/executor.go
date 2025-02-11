package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var excStr string
	for _, cm := range cmd {
		excStr += fmt.Sprintf(" %s", cm)
	}
	excStr = strings.TrimLeft(excStr, " ")
	slog.Info("Строка запуска", "cmd", excStr)
	exc := exec.Command(excStr)
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
