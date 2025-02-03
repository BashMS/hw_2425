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

func checkSrcFile(fromPath string) error {
	fileStt, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("os.Stat: %w", err)
	}

	// проверим размер и возможность чтения
	if fileStt.Size() == 0 {
		file, err := os.Open(fromPath)
		if err != nil {
			return fmt.Errorf("os.Open: %w", err)
		}
		defer file.Close()
		buf := make([]byte, 1)
		n, err := file.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("file.Read: %w", err)
		}
		// Если что то прочитали из нулевого файла, то он какой то не такой
		if n > 0 {
			return ErrUnsupportedFile
		}
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := checkSrcFile(fromPath)
	if err != nil {
		return fmt.Errorf("checkSrcFile: %w", err)
	}

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
