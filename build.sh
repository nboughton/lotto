#!/bin/bash
## Rebuild existing code
echo "Rebuilding all code"
go build -o site.app
if [ $? -ne 0 ]; then
  echo "Build error!"
  exit 1
fi

cd frontend
npm run build
if [ $? -ne 0 ]; then
  echo "Build error!"
  exit 1
fi

cd ..