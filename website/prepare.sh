#!/bin/sh
set -xe

cp ../README.md content/_index.md
cp ../cmd/cs/README.md content/cs/_index.md
mkdir -p static
cp ../gopher.png static

cd ../
dep ensure -vendor-only

cd cmd/cs
dep ensure -vendor-only
go build
./cs gen-doc

set +xe
echo "we are now ready to: hugo serve"
