#!/usr/bin/env sh
rootdir=$(pwd)

cd gen/consumerapi
rm *Get.go

# generate the consumer API
cd $rootdir
cd cmd/generate/
go build .
./generate

# fix imports
cd $rootdir
cd gen/consumerapi
goimports -w *Get.go

# run the example
cd $rootdir
cd cmd/main/
go build .
./main
