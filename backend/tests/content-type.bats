
function cleanup {
  unlink out || true
  curl --max-time 5 -X DELETE http://backend:8080/file.txt
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

@test "automatic content type txt" {
  run curl --fail --max-time 5 -X PUT -T file.txt http://backend:8080/file.txt
  [ "$status" -eq 0 ]

  run sh -c 'curl --fail --max-time 5 --silent -v http://backend:8080/file.txt 2>out'
  [ "$status" -eq 0 ]

  run grep "Content-Type: text/plain" out
  [ "$status" -eq 0 ]
}

@test "automatic content type png" {
  run curl --fail --max-time 5 -X PUT -T favicon.png http://backend:8080/favicon.png
  [ "$status" -eq 0 ]

  run sh -c 'curl --fail --max-time 5 --silent -v http://backend:8080/favicon.png 2>out'
  [ "$status" -eq 0 ]

  run grep "Content-Type: image/png" out
  [ "$status" -eq 0 ]
}

@test "manual content type" {
  run curl --fail --max-time 5 -X PUT -H "Content-Type: dinghy/test" -T file.txt http://backend:8080/file.txt
  [ "$status" -eq 0 ]

  run sh -c 'curl --fail --max-time 5 --silent -v http://backend:8080/file.txt 2>out'
  [ "$status" -eq 0 ]

  run grep "Content-Type: dinghy/test" out
  [ "$status" -eq 0 ]
}

@test "automatic content type txt presigned" {
  run curl --fail --max-time 5 -X PUT -T file.txt http://backend:8080/file.txt
  [ "$status" -eq 0 ]

  run sh -c 'curl -L --fail --max-time 5 --silent -v "http://backend:8080/file.txt?redirect" 2>out'
  [ "$status" -eq 0 ]

  run grep "Content-Type: text/plain" out
  [ "$status" -eq 0 ]
}

@test "automatic content type png presigned" {
  run curl --fail --max-time 5 -X PUT -T favicon.png http://backend:8080/favicon.png
  [ "$status" -eq 0 ]

  run sh -c 'curl -L --fail --max-time 5 --silent -v "http://backend:8080/favicon.png?redirect" 2>out'
  [ "$status" -eq 0 ]

  run grep "Content-Type: image/png" out
  [ "$status" -eq 0 ]
}

@test "manual content type presigned" {
  run curl --fail --max-time 5 -X PUT -H "Content-Type: dinghy/test" -T file.txt http://backend:8080/file.txt
  [ "$status" -eq 0 ]

  run sh -c 'curl -L --fail --max-time 5 --silent -v "http://backend:8080/file.txt?redirect" 2>out'
  [ "$status" -eq 0 ]

  run grep "Content-Type: dinghy/test" out
  [ "$status" -eq 0 ]
}
