
function cleanup {
  curl --max-time 5 -X DELETE http://backend:8080/favicon.png
  unlink download.png || true
}

setup () {
  curl --silent --fail --max-time 5 -X PUT -T favicon.png http://backend:8080/favicon.png
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

@test "redirect GET" {
  run curl -v --silent --fail --max-time 5 -o download.png "http://backend:8080/favicon.png?redirect"
  [ "$status" -eq 0 ]

  run sh -c "echo \"$output\" | grep \"307 Temporary Redirect\""
  [ "$status" -eq 0 ]
}

@test "redirect PUT" {
  run curl -v --silent --fail --max-time 5 -X PUT -T favicon.png "http://backend:8080/favicon.png?redirect"
  [ "$status" -eq 0 ]

  run sh -c "echo \"$output\" | grep \"307 Temporary Redirect\""
  [ "$status" -eq 0 ]
}

@test "redirect DELETE" {
  run curl -v --silent --fail --max-time 5 -X DELETE "http://backend:8080/favicon.png?redirect"
  [ "$status" -eq 0 ]

  run sh -c "echo \"$output\" | grep \"307 Temporary Redirect\""
  [ "$status" -eq 0 ]
}
