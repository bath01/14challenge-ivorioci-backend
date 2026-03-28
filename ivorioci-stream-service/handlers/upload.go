package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	maxVideoSize = 2 << 30  // 2 GB
	maxThumbSize = 10 << 20 // 10 MB
)

var allowedVideoTypes = map[string]string{
	"video/mp4":       ".mp4",
	"video/webm":      ".webm",
	"video/ogg":       ".ogv",
	"video/quicktime": ".mov",
	"video/x-msvideo": ".avi",
}

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func detectMIME(f multipart.File) (string, error) {
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	mimeType := http.DetectContentType(buf[:n])
	_, err = f.Seek(0, io.SeekStart)
	return mimeType, err
}

func saveUploadedFile(f multipart.File, dst string) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(dst), 0o750); err != nil {
		return 0, fmt.Errorf("mkdir: %w", err)
	}
	out, err := os.Create(dst)
	if err != nil {
		return 0, fmt.Errorf("create: %w", err)
	}
	defer out.Close()
	return io.Copy(out, f)
}

func uniqueFilename(ext string) string {
	return uuid.New().String() + ext
}
