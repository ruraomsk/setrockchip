#!/bin/bash
echo "Start to Windows installer"
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.buildDate=$(date -u +%Y-%m-%d)"
if [ $? -ne 0 ]; then
	echo 'An error has occurred! Aborting the script execution...'
	exit 1
fi
echo "Start to Linux installer"
go build -ldflags "-X main.buildDate=$(date -u +%Y-%m-%d)"
if [ $? -ne 0 ]; then
	echo 'An error has occurred! Aborting the script execution...'
	exit 1
fi
