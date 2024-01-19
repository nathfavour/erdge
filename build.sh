#!/bin/bash

# Delete the binaries if they already exist
rm -f erdgeWin erdgeMac erdgeLin

# Compile the main.go file for different platforms
# For Windows
GOOS=windows GOARCH=amd64 go build -o erdgeWin main.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o erdgeMac main.go

# For Linux
GOOS=linux GOARCH=amd64 go build -o erdgeLin main.go

echo "Compilation completed."