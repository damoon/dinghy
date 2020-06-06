
teardown () {
  echo teardown log
  echo "exit code: $status"
  for i in "${!lines[@]}"; do 
    printf "line %s:\t%s\n" "$i" "${lines[$i]}"
  done
  echo teardown done
}

@test "redirect" {
  run curl -v --silent --fail --max-time 5 --user-agent "pretend-to-be-a-browser" http://backend/
  [ "$status" -eq 0 ]

  run sh -c "echo \"$output\" | grep \"307 Temporary Redirect\""
  [ "$status" -eq 0 ]
}
