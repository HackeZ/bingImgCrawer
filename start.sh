#!/bin/bash
## Author: HackerZ

## FOR Golang Lib.
go get gopkg.in/alecthomas/kingpin.v2

## Delete Older Source.
if [ -f "./bingImgCrawer" ];then
    rm "bingImgCrawer"
fi

## Build Application.
go build -o bingImgCrawer main.go

## Run.
if [ -f "./bingImgCrawer" ];then
    ./bingImgCrawer --root "./"
fi