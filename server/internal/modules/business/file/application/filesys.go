package application

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func sanitizeFilename(original string) string {
	ext := filepath.Ext(original)
	base := strings.TrimSuffix(original, ext)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(base)))
	return fmt.Sprintf("%s_%s%s", hash[:8], time.Now().Format("20060102150405"), ext)
}

func mkdirAll(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func writeBytesToFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0666)
}

func remove(path string) error {
	return os.Remove(path)
}
