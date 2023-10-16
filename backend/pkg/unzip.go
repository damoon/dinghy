package dinghy

import (
	"context"
	"fmt"
	"os"
	"strings"

	archiver "github.com/mholt/archiver/v3"
	"github.com/opentracing/opentracing-go"
)

func (s ServiceServer) unzip(ctx context.Context, path string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "extract object")
	defer span.Finish()

	ext, err := archiveExtension(path)
	if err != nil {
		return err
	}

	target := filesDirectory + strings.TrimSuffix(path, ext)

	exists, _, _, err := s.Storage.exists(ctx, target)
	if err != nil {
		return fmt.Errorf("verify target location %s: %v", target, err)
	}

	if exists {
		return fmt.Errorf("target %s exists", target)
	}

	tmpfile, err := os.CreateTemp("", "s3_download_*"+ext)
	if err != nil {
		return fmt.Errorf("create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	err = s.Storage.download(ctx, filesDirectory+path, tmpfile)
	if err != nil {
		return fmt.Errorf("download: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "s3_download")
	if err != nil {
		return fmt.Errorf("create temp file: %v", err)
	}
	defer os.Remove(tmpDir)

	{
		span, _ := opentracing.StartSpanFromContext(ctx, "extract archive")
		defer span.Finish()

		err = archiver.Unarchive(tmpfile.Name(), tmpDir)
		if err != nil {
			return fmt.Errorf("extract file: %v", err)
		}
	}

	defer span.Finish()
	err = s.Storage.uploadRecursive(ctx, tmpDir, target)
	if err != nil {
		return fmt.Errorf("upload: %v", err)
	}

	return nil
}

func canBeExtracted(file string, dirs []string) bool {
	ext, err := archiveExtension(file)
	if err != nil {
		return false
	}

	base := strings.TrimSuffix(file, ext)

	for _, dir := range dirs {
		if base == dir {
			return false
		}
	}

	return true
}

func archiveExtension(path string) (string, error) {
	extensions := []string{
		".rar",
		".tar",
		".tar.br",
		".tbr",
		".tar.bz2",
		".tbz2",
		".tar.gz",
		".tgz",
		".tar.lz4",
		".tlz4",
		".tar.sz",
		".tsz",
		".tar.xz",
		".txz",
		".tar.zst",
		".zip",
	}

	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return ext, nil
		}
	}

	return "", fmt.Errorf("archive type not supported")
}
