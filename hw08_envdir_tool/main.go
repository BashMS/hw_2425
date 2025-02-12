package main

import (
	"fmt"
	"log/slog"
	"os"
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

func main() {
	slog.Info("Start...")
	args := os.Args
	if len(args) == 1 {
		slog.Info("Finish.")
		return
	}
	// Соберем переменные
	env, err := ReadDir(args[1])
	if err != nil {
		slog.Error("ReadDir: ", "Error", err)
		return
	}

	if len(args) > 2 {
		err = checkVariables(env)
		if err != nil {
			slog.Error("Ошибка обработки переменных окружения", "checkVariables", err)
		}
		code := RunCmd(args[2:], env)
		slog.Info("Завершили", "ReturnCode", code)
	}
}
