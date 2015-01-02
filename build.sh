#!/bin/bash

printf "** Building linux/amd64\n"
go build -a -o bin/linux-amd64/up53 github.com/zerklabs/up53
