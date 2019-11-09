#!/usr/bin/env bats

teardown () {
  echo teardown log
  echo "exit code: $status"
  for i in "${!lines[@]}"; do 
    printf "line %s:\t%s\n" "$i" "${lines[$i]}"
  done
  echo teardown done
}

@test "upload and download" {
  run curl --fail -L -T tests/test.txt http://127.0.0.1:8080/upload-test
  [ $status -eq 0 ]

  run curl --fail -L -o /tmp/upload-download-test http://127.0.0.1:8080/upload-test
  [ $status -eq 0 ]

  run diff tests/test.txt /tmp/upload-download-test
  [ $status -eq 0 ]
}
