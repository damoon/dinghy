#!/usr/bin/env bats

load helpers

@test "redirect to minio: root" {
  run sh -c "curl --fail -v http://127.0.0.1:8080/ 2>&1 | grep 'HTTP/1.1 307 Temporary Redirect'"
  [ $status -eq 0 ]
}

@test "redirect to minio: file" {
  run sh -c "curl --fail -v http://127.0.0.1:8080/some-not-existing-file 2>&1 | grep 'HTTP/1.1 307 Temporary Redirect'"
  [ $status -eq 0 ]
}

@test "redirect to minio: path" {
  run sh -c "curl --fail -v http://127.0.0.1:8080/some-not/existing-file 2>&1 | grep 'HTTP/1.1 307 Temporary Redirect'"
  [ $status -eq 0 ]
}

@test "redirect to minio: directory" {
  run sh -c "curl --fail -v http://127.0.0.1:8080/some-dir/ 2>&1 | grep 'HTTP/1.1 307 Temporary Redirect'"
  [ $status -eq 0 ]
}
