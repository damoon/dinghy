
function cleanup {
  curl --max-time 5 -X DELETE http://backend/favicon.png
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

@test "list" {
  run curl --fail --max-time 5 -o list http://backend/
  [ "$status" -eq 0 ]

  run grep favicon.png list
  [ "$status" -ne 0 ]

  run curl --fail --max-time 5 -X PUT -T favicon.png http://backend/favicon.png
  [ "$status" -eq 0 ]

  run curl --fail --max-time 5 -o list http://backend/
  [ "$status" -eq 0 ]

  run grep favicon.png list
  [ "$status" -eq 0 ]

  run curl --fail --max-time 5 -X DELETE http://backend/favicon.png
  [ "$status" -eq 0 ]

  run curl --fail --max-time 5 -o list http://backend/
  [ "$status" -eq 0 ]

  run grep favicon.png list
  [ "$status" -ne 0 ]
}
