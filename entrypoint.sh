#!/bin/sh
until /init
do
  echo "Init failed... waiting 5 seconds before retry"
  sleep 5
done
