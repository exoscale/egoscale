#!/bin/sh
set -xe

cp ../README.md content/_index.md
mkdir -p static
cp ../gopher.png static

set +xe
echo "we are now ready to: hugo serve"
