package main

import (
	"log/slog"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		return
	}
	// Соберем переменные
	env, err := ReadDir(args[1])
	if err != nil {
		slog.Error("ReadDir: ", "Error", err)
		return
	}

	if len(args) > 2 {
		if err != nil {
			slog.Error("Ошибка обработки переменных окружения", "checkVariables", err)
		}
		code := RunCmd(args[2:], env)
		slog.Info("Завершили", "ReturnCode", code)
		os.Exit(code)
	}
}
