#!/usr/bin/env bats

load helpers

@test "upload and download" {
  run curl --fail -L -T tests/helpers.bash http://127.0.0.1:8080/upload-test
  [ $status -eq 0 ]

  run curl --fail -L -o /tmp/upload-download-test http://127.0.0.1:8080/upload-test
  [ $status -eq 0 ]

  run diff tests/helpers.bash /tmp/upload-download-test
  [ $status -eq 0 ]
}
