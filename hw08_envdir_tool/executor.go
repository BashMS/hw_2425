package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// checkVariables проверка и установка переменных окружения.
func checkVariables(env Environment) error {
	var err error
	for k, v := range env {
		// удалим если нужно
		if v.NeedRemove {
			err = os.Unsetenv(k)
			if err != nil {
				return fmt.Errorf("os.Unsetenv: %w", err)
			}
			continue
		}
		// иначе пересоздаем
		err = os.Unsetenv(k)
		if err != nil {
			return fmt.Errorf("os.Unsetenv: %w", err)
		}
		err = os.Setenv(k, v.Value)
		if err != nil {
			return fmt.Errorf("os.Setenv: %w", err)
		}
	}
	return nil
}

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
	exc.Stderr = os.Stderr
	err = exc.Run()
	if err != nil {
		var exiterr *exec.ExitError
		if errors.As(err, &exiterr) {
			returnCode = exiterr.ExitCode()
		} else {
			slog.Error("exc.Run", "Error", err)
			returnCode = 1
		}
	}

	return
}
