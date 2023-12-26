#!/bin/sh
sigint_handler()
{
  pkill -TERM -P $PID
  exit
}

trap sigint_handler SIGINT

while true; do
  go generate ./...
  go run cmd/main.go &
  PID=$!
  sleep 3
  fswatch -r1 -e "\\.swp$" -e "\\.swx$" .
  pkill -TERM -P $PID
done

