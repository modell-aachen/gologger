#!/bin/sh

rootPwd=$(pwd)

go build
cd $rootPwd/cmd/logstore && go build && \
cd $rootPwd/cmd/logreport && go build && \
cd $rootPwd/cmd/logtail && go build && \
cd $rootPwd/cmd/logdump && go build

