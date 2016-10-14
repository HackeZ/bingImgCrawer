#!/bin/bash
## Author: HackerZ

## FOR Golang Lib
go get gopkg.in/alecthomas/kingpin.v2

## Build Application
go build -o bingImgCrawer main.go

## Run
./bingImgCrawer -root "./"