package helper

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const MaxImageSize = 2 << 20

var (
	ErrImageTooLarge   = errors.New("ukuran gambar melebihi batas maksimal 2MB")
	ErrInvalidImageType = errors.New("format gambar harus .jpg atau .png")
)

func imageExt(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}

func ValidateImage(fh *multipart.FileHeader) error {
	if fh.Size > MaxImageSize {
		return ErrImageTooLarge
	}

	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := io.ReadFull(f, buf)
	if imageExt(http.DetectContentType(buf[:n])) == "" {
		return ErrInvalidImageType
	}
	return nil
}

// SaveImage menyimpan file gambar ke dir dengan nama filename + ekstensi
// hasil deteksi konten, lalu mengembalikan path relatif (slash).
func SaveImage(fh *multipart.FileHeader, dir, filename string) (string, error) {
	if err := ValidateImage(fh); err != nil {
		return "", err
	}

	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	buf := make([]byte, 512)
	n, _ := io.ReadFull(f, buf)
	ext := imageExt(http.DetectContentType(buf[:n]))

	path := filepath.Join(dir, filename+ext)
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := dst.Write(buf[:n]); err != nil {
		return "", err
	}
	if _, err := io.Copy(dst, f); err != nil {
		return "", err
	}
	return filepath.ToSlash(path), nil
}

func RemoveFile(path string) {
	if path == "" {
		return
	}
	clean := filepath.Clean(strings.TrimPrefix(path, "/"))
	_ = os.Remove(clean)
}
