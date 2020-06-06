
function cleanup {
  curl --max-time 5 -X DELETE http://backend:8080/favicon.png
  unlink download.png || true
}

setup () {
    cleanup
}

teardown () {
  cleanup

  echo teardown log
  echo "exit code: $status"
  for i in "${!lines[@]}"; do 
    printf "line %s:\t%s\n" "$i" "${lines[$i]}"
  done
  echo teardown done
}

@test "put get delete" {
  run curl --fail --max-time 5 -o download.png http://backend:8080/favicon.png
  [ "$status" -eq 0 ]

  run diff favicon.png download.png
  [ "$status" -ne 0 ]

  run curl --fail --max-time 5 -X PUT -T favicon.png http://backend:8080/favicon.png
  [ "$status" -eq 0 ]

  run curl --fail --max-time 5 -o download.png http://backend:8080/favicon.png
  [ "$status" -eq 0 ]

  run diff favicon.png download.png
  [ "$status" -eq 0 ]

  run curl --fail --max-time 5 -X DELETE http://backend:8080/favicon.png
  [ "$status" -eq 0 ]

  run curl --fail --max-time 5 -o download.png http://backend:8080/favicon.png
  [ "$status" -eq 0 ]

  run diff favicon.png download.png
  [ "$status" -ne 0 ]
}
