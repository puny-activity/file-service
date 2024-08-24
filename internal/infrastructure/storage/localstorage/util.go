package localstorage

import (
	"crypto/md5"
	"fmt"
	"github.com/h2non/filetype"
	"io"
	"os"
)

func calculateMD5(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

func detectContentType(file *os.File) (string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return "", err
	}

	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		return "application/octet-stream", nil
	}
	return kind.MIME.Value, nil
}
