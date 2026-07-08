#!/bin/bash
export Result=1;
while [ $Result -ne 0 ]; do
  echo "STARTING ($Result) @" `date`;
  rsync -Pza --timeout=10 "$@";
  Result=$?;
  sleep 1;
done