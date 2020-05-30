package dinghy

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

var (
	s3Endpoint  = flag.String("s3Endpoint", "", "run database integration tests")
	s3AccessID  = flag.String("s3AccessID", "", "run message queue integration tests")
	s3AccessKey = flag.String("s3AccessKey", "", "run message queue integration tests")
	s3UseSSL    = flag.Bool("s3UseSSL", true, "run s3")
	s3Bucket    = flag.String("s3Bucket", "", "run message queue integration tests")
	s3Location  = flag.String("s3Location", "us-east-1", "run message queue integration tests")
)

func TestSignedUploadDownload(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	flag.Parse()

	if *s3Endpoint == "" {
		t.Errorf("-s3Endpoint not set")
	}
	if *s3AccessID == "" {
		t.Errorf("-s3AccessID not set")
	}
	if *s3AccessKey == "" {
		t.Errorf("-s3AccessKey not set")
	}
	if *s3Bucket == "" {
		t.Errorf("-s3Bucket not set")
	}

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	minioClient, err := minio.New(*s3Endpoint, *s3AccessID, *s3AccessKey, *s3UseSSL)
	if err != nil {
		t.Fatalf("creating minio client: %v", err)
	}
	minioClient.SetCustomTransport(transport)

	var s presignStorage = NewMinioStorage(minioClient, *s3Bucket, *s3Location)

	found, err := minioClient.BucketExists(*s3Bucket)
	if err != nil {
		t.Fatalf("checking bucket: %v", err)
	}
	if !found {
		err = minioClient.MakeBucket(*s3Bucket, *s3Location)
		if err != nil {
			t.Fatalf("creating bucket: %v", err)
		}
	}

	// ensure free filename
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	name := uuid.New()
	found, err = s.exists(ctx, name.String())
	if err != nil {
		t.Fatalf("s.exists() failed: %v", err)
	}
	if found {
		t.Fatalf("s.exists() = %v, want %v", found, false)
	}

	// put file
	url, err := s.presign(http.MethodPut, name.String())
	if err != nil {
		t.Fatalf("presign PUT: %v", err)
	}

	testFileContent := []byte("Hello, world.")
	putReq, err := http.NewRequest(http.MethodPut, url.String(), bytes.NewReader(testFileContent))
	if err != nil {
		t.Fatalf("initiale PUT request: %v", err)
	}
	putReq.Header.Add("Content-Length", strconv.Itoa(len(testFileContent)))

	httpClient := &http.Client{
		Timeout:   time.Second * 20,
		Transport: transport,
	}
	resp, err := httpClient.Do(putReq)
	if err != nil {
		t.Fatalf("PUT request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT request status code %v, want %v", resp.StatusCode, http.StatusOK)
	}

	// file exists
	found, err = s.exists(ctx, name.String())
	if err != nil {
		t.Fatalf("s.exists() failed: %v", err)
	}
	if !found {
		t.Fatalf("s.exists() = %v, want %v", found, true)
	}

	// download and compare
	url, err = s.presign(http.MethodGet, name.String())
	if err != nil {
		t.Fatalf("presign GET: %v", err)
	}
	resp, err = http.Get(url.String())
	if err != nil {
		t.Fatalf("GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET request status code %v, want %v", resp.StatusCode, http.StatusOK)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET request status code %v, want %v", resp.StatusCode, http.StatusOK)
	}
	returnedFile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read GET request: %v", err)
	}
	if !reflect.DeepEqual(returnedFile, testFileContent) {
		t.Fatalf("returned file %v, want %v", returnedFile, testFileContent)
	}

	// cleanup
	err = minioClient.RemoveObject(*s3Bucket, name.String())
	if err != nil {
		t.Fatalf("cleanup: remove test file: %v", err)
	}
}
