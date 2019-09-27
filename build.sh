#!/bin/bash
## Rebuild existing code
echo "Building quasar frontend: public/"
cd frontend
quasar b
if [ $? -ne 0 ]; then
  echo "Build error!"
  exit 1
fi

cd ..

echo "Building Go binary: site.app"
go build -o site.app
if [ $? -ne 0 ]; then
  echo "Build error!"
  exit 1
fi

if [[ $1 == "-d" ]]; then
  docker build . -t nboughton/lotto --no-cache
fi