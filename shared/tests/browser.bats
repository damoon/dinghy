
setup () {
  curl --silent --fail --max-time 5 -X PUT -T favicon.png http://dinghy-backend:8080/subfolder/favicon.png
}

teardown () {
  curl --silent --fail --max-time 5 -X DELETE http://dinghy-backend:8080/subfolder/favicon.png

  echo teardown log
  echo "exit code: $status"
  for i in "${!lines[@]}"; do 
    printf "line %s:\t%s\n" "$i" "${lines[$i]}"
  done
  echo teardown done
}

@test "from backend redirect to uploaded file in folder" {
  run main
  [ "$status" -eq 0 ]
}
