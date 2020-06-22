package dinghy

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	archiver "github.com/mholt/archiver/v3"
)

func (s ServiceServer) unzip(path string) error {
	ext, err := archiveExtension(path)
	if err != nil {
		return err
	}

	target := filesDirectory + strings.TrimSuffix(path, ext)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	exists, _, _, err := s.Storage.exists(ctx, target)
	if err != nil {
		return fmt.Errorf("verify target location %s: %v", target, err)
	}

	if exists {
		return fmt.Errorf("target %s exists", target)
	}

	tmpfile, err := ioutil.TempFile("", "s3_download_*"+ext)
	if err != nil {
		return fmt.Errorf("create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	err = s.Storage.download(ctx, filesDirectory+path, tmpfile)
	if err != nil {
		return fmt.Errorf("download: %v", err)
	}

	tmpDir, err := ioutil.TempDir("", "s3_download")
	if err != nil {
		return fmt.Errorf("create temp file: %v", err)
	}
	defer os.Remove(tmpDir)

	err = archiver.Unarchive(tmpfile.Name(), tmpDir)
	if err != nil {
		return fmt.Errorf("extract file: %v", err)
	}

	err = s.Storage.uploadRecursive(ctx, tmpDir, target)
	if err != nil {
		return fmt.Errorf("upload: %v", err)
	}

	return nil
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
