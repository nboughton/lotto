#!/bin/bash

## Rebuild existing code
echo "Rebuilding"
go build -v -o site.app
if [ $? -ne 0 ]; then
  echo "BUILD ERROR"
  exit 1
fi

cd frontend
npm run build
if [ $? -ne 0 ]; then
  echo "BUILD ERROR"
  exit 1
fi
cd ..
