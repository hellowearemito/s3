#!/bin/sh

set -e
if [ -z "$PKG_NAME" ]; then
	echo "Package name (\$PKG_NAME) is not set. Aborting."
	exit 1
fi

set -x

mkdir -p "$(dirname ".go/src/$PKG_NAME")"
rm -f ".go/src/$PKG_NAME"
ln -s "$(readlink -f .)" ".go/src/$PKG_NAME"
export GOPATH
GOPATH=$(readlink -f .)/.go
PATH="$GOPATH/bin:$PATH"
cd "$GOPATH/src/$PKG_NAME" || exit 1
go get -u github.com/golang/dep/cmd/dep
go get -u github.com/mibk/g..
dep ensure