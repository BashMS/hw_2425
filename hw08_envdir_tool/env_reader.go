package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func getValue(filePath string) (string, error) {
	fBody, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("os.Open: %w", err)
	}
	defer fBody.Close()

	var result string
	rBuffer := make([]byte, 1)
	for {
		_, err := fBody.Read(rBuffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return "", fmt.Errorf("reader.ReadString: %w", err)
		}
		if rBuffer[0] != '\n' && rBuffer[0] != '\r' && err != io.EOF {
			result += string(rBuffer)
			continue
		}
		break
	}

	result = strings.TrimRight(result, " ") // пробел
	result = strings.TrimRight(result, "	") // таб

	i := strings.Index(result, "\x00")
	if i != -1 {
		resultC := fmt.Sprintf("%s\n%s", result[:i], result[i+1:])
		result = resultC
	}
	return result, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Прочитаем указанную директорию
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %w", err)
	}
	// Сформируем мапу переменных окружения
	result := make(Environment, len(files))
	for _, file := range files {
		fInfo, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("file.Info(): %w", err)
		}
		// Если пустой тогда пометим на удаление
		if fInfo.Size() == 0 {
			result[file.Name()] = EnvValue{NeedRemove: true}
			continue
		}
		val, err := getValue(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("getValue: %w", err)
		}
		if len(val) == 0 {
			result[file.Name()] = EnvValue{
				Value:      val,
				NeedRemove: true,
			}
			continue
		}
		result[file.Name()] = EnvValue{
			Value:      val,
			NeedRemove: false,
		}
	}
	return result, nil
}
