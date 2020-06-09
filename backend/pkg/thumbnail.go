package dinghy

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"  // load gif image support
	_ "image/jpeg" // load jpeg image support
	"image/png"    // load png image support
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/opentracing/opentracing-go"
	_ "golang.org/x/image/bmp"  // load bmp image support
	_ "golang.org/x/image/tiff" // load tiff image support
	_ "golang.org/x/image/webp" // load webp image support
)

const thumbnailWidth = 122
const thumbnailHeight = 72

const thumbnailDirectory = "thumbnails"
const filesDirectory = "files"

func thumbnailSupported(path string) bool {
	ext := filepath.Ext(path)

	switch ext {
	case ".gif":
	case ".jpg":
	case ".jpeg":
	case ".png":
	case ".bmp":
	case ".tiff":
	case ".webp":
	default:
		return false
	}

	return true
}

func (s *ServiceServer) prepareThumbnail(ctx context.Context, etag, path string) (string, error) {
	if !thumbnailSupported(path) {
		return "", fmt.Errorf("extention not supported")
	}

	thumbnailPath := thumbnailDirectory + "/" + etag + ".png"

	exists, _, err := s.Storage.exists(ctx, thumbnailPath)
	if err != nil {
		return "", fmt.Errorf("lookup cached thumbnail: %v", err)
	}

	if exists {
		return thumbnailPath, nil
	}

	tmpfile, err := ioutil.TempFile("", "s3_download")
	if err != nil {
		return "", fmt.Errorf("create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	err = s.Storage.download(ctx, filesDirectory+path, tmpfile)
	if err != nil {
		return "", fmt.Errorf("download: %v", err)
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("seek file: %v", err)
	}

	err = resizeImage(ctx, tmpfile)
	if err != nil {
		return "", fmt.Errorf("resize image: %v", err)
	}

	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("seek thumbnail temp file: %v", err)
	}

	err = s.Storage.upload(ctx, thumbnailPath, tmpfile)
	if err != nil {
		return "", fmt.Errorf("upload thumbnail: %v", err)
	}

	return thumbnailPath, nil
}

func resizeImage(ctx context.Context, f *os.File) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "create thumbnail")
	defer span.Finish()

	in, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	spanCompute, ctx := opentracing.StartSpanFromContext(ctx, "compute thumbnail")
	out := imaging.Fit(in, thumbnailWidth, thumbnailHeight, imaging.Lanczos)
	spanCompute.Finish()

	_, err = f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seek thumbnail file: %v", err)
	}

	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("truncate thumbnail file: %v", err)
	}

	err = png.Encode(f, out)
	if err != nil {
		return fmt.Errorf("save thumbnail back to file: %v", err)
	}

	return nil
}
