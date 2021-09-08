#!/usr/bin/env bash

i=0
while true
do
  now=$(date +"%Y-%m-%d %T")
  printf '{"i": "%s", "time": "%s","level":"info"}\n' "$i" "$now"
  sleep 1
  ((i=i+1))
done