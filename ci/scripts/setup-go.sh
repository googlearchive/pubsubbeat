#!/bin/bash

PROGNAME=pubsubbeat
GODIR=github.com/GoogleCloudPlatform/pubsubbeat

mkdir -p $GOPATH/src/github.com/GoogleCloudPlatform
ln -s $PWD/src $GOPATH/src/$GODIR

cd $GOPATH/src/$GODIR
echo "Gopath is: " $GOPATH
echo "pwd is: " $PWD
ls -lah

