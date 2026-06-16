package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveBase64Image(base64Data, prefix string) (string, error) {
	if base64Data == "" {
		return "", nil
	}

	// contoh: data:image/jpeg;base64,xxxx
	parts := strings.Split(base64Data, ",")
	if len(parts) != 2 {
		return "", errors.New("invalid base64 image format")
	}

	// ambil mime type
	header := parts[0] // data:image/jpeg;base64
	headerParts := strings.Split(header, ";")
	if len(headerParts) == 0 {
		return "", errors.New("invalid base64 header")
	}

	mimeType := strings.TrimPrefix(headerParts[0], "data:")

	// mapping mime → extension
	var ext string
	switch mimeType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	case "application/pdf":
		ext = ".pdf"
	default:
		return "", errors.New("unsupported image type")
	}

	// decode data
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// pastikan folder ada
	dir := "static/uploads/absensi"
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating directory", err)
		return "", err
	}

	// generate nama file
	filename := prefix + "_" + time.Now().Format("20060102_150405") + ext
	fullpath := filepath.Join(dir, filename)

	fmt.Println("Saving image to", fullpath)

	if err := os.WriteFile(fullpath, imageData, 0644); err != nil {
		fmt.Println("Error write to file", err)
		return "", err
	}

	return fullpath, nil
}

func SaveBase64ImageInv(base64Data, prefix string) (string, error) {
	if base64Data == "" {
		return "", nil
	}

	// contoh: data:image/jpeg;base64,xxxx
	parts := strings.Split(base64Data, ",")
	if len(parts) != 2 {
		return "", errors.New("invalid base64 image format")
	}

	// ambil mime type
	header := parts[0] // data:image/jpeg;base64
	headerParts := strings.Split(header, ";")
	if len(headerParts) == 0 {
		return "", errors.New("invalid base64 header")
	}

	mimeType := strings.TrimPrefix(headerParts[0], "data:")

	// mapping mime → extension
	var ext string
	switch mimeType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	case "application/pdf":
		ext = ".pdf"
	default:
		return "", errors.New("unsupported image type")
	}

	// decode data
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// pastikan folder ada
	dir := "static/uploads/inventory"
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating directory", err)
		return "", err
	}

	// generate nama file
	filename := prefix + "_" + time.Now().Format("20060102_150405") + ext
	fullpath := filepath.Join(dir, filename)

	fmt.Println("Saving image to", fullpath)

	if err := os.WriteFile(fullpath, imageData, 0644); err != nil {
		fmt.Println("Error write to file", err)
		return "", err
	}

	return fullpath, nil
}
