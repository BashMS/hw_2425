package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Получим размер копируемого файла
	fileStt, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("os.Stat: %w", err)
	}
	// Если указали отступ тогда проверим его валидность
	if offset > fileStt.Size() {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || limit > fileStt.Size() {
		limit = fileStt.Size()
	}

	// Открываем файл
	fileSrc, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer fileSrc.Close()
	// устанавливаем смещение
	if offset > 0 {
		_, err = fileSrc.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("file.Seek: %w", err)
		}
	}
	// Создаем файл приемник
	fileDst, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer fileDst.Close()

	data := make([]byte, 1)

	progressBar := pb.Start64(limit)
	defer progressBar.Finish()
	var curCnt int64
	for {
		_, err := fileSrc.Read(data)
		if errors.Is(err, io.EOF) { // если конец файла
			break // выходим из цикла
		}
		wb, err := fileDst.Write(data)
		if err != nil {
			return fmt.Errorf("fileDst.Write: %w", err)
		}
		curCnt += int64(wb)
		progressBar.Add(wb)
		if curCnt >= limit {
			break
		}
	}

	return nil
}
