
setup () {
  curl --silent --fail --max-time 5 -X PUT -T hello-world.txt http://backend:8080/subfolder/hello-world.txt
}

teardown () {
  curl --silent --fail --max-time 5 -X DELETE http://backend:8080/subfolder/hello-world.txt

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
