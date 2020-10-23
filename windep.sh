#!/bin/bash
echo "Start to Windows deploy"
GOOS=windows GOARCH=amd64 go build
cp *.exe /home/rura/vm/windows
