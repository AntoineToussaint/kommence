#!/usr/bin/env bash

i=0
while true
do
  now=$(date +"%Y-%m-%d %T")
  printf '{"i": "%s", "time": "%s","level":"info", "env": "%s"\n' "$i" "$now" "${KOMMENCE_VAR}"
  sleep 1
  ((i=i+1))
done