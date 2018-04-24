#!/bin/sh

export TERM=xterm-color

. ./ci-prepare.sh

set -e
set -x

g.. build
g.. install

set +e

if fmtresults=$(find . -type d \( -path "*/vendor" -o -path "*/.*" \) -prune -o -iname "*.go" -print0 | xargs -I{} -0 gofmt -l '{}'); then
	# $() removes trailing newlines. Reinsert it if the string is not empty
	# so that wc -l can count the last line.
	if [ -n "$fmtresults" ]; then
		newline=$(printf "\\nX") # protect the newline from $() using a trailing X
		newline=${newline%X} # remove the X
		fmtresults="$fmtresults$newline"
	fi
	efmt=$(printf "%s" "$fmtresults" | wc -l)
else
	efmt=1
fi

g.. vet
evet=$?
# g.. test
rm -f coverage.tmp
# shellcheck disable=SC2016
echo 'mode: set' > coverage.txt && g.. | xargs -n1 -I{} sh -c 'go test -v -coverprofile=coverage.tmp {}; ec=$?; if test -f coverage.tmp; then tail -n +2 coverage.tmp >> coverage.txt; fi; exit $ec'
etest=$?
rm coverage.tmp
go tool cover -func coverage.txt | tail -1

set +x
if [ "$efmt" != "0" ]; then
	echo "Code is not formatted correctly. Please run go fmt."
fi
if [ "$efmt" != "0" ] || [ $evet -ne 0 ] || [ $etest -ne 0 ]; then
	exit 1
fi