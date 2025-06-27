#script to build file in linux & windows

#!/bin/bash

go mod init te

export GOOS=linux
go build -o te te.go

export GOOS=windows
go build -o te.exe te.go