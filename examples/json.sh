#!/usr/bin/env bash

i=0
while true
do
  echo \{\"j\": $i\}
  sleep 1
  ((i=i+1))
done