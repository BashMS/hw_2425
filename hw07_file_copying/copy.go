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
	ErrSpecifiedFilesMatches = errors.New("specified files matches")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Получим размер копируемого файла
	fileStt, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("os.Stat: %w", err)
	}

	// Проверим файл приемник
	fileDstStt, err := os.Stat(toPath)
	switch {
	case err == nil:
		if os.SameFile(fileStt, fileDstStt) {
			return ErrSpecifiedFilesMatches
		}
	case !errors.Is(err, os.ErrNotExist):
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

	bar := pb.Full.Start64(limit)
	defer bar.Finish()

	barReader := bar.NewProxyReader(fileSrc)
	_, err = io.CopyN(fileDst, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) { // если конец файла
		return fmt.Errorf("io.CopyN: %w", err)
	}

	return nil
}
