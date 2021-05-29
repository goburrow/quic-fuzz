#!/bin/sh
set -e
if [ $# -ne 1 ]; then
    echo "$0" "[ClientInitial,ServerInitial,Client,Server]"
    exit 1
fi
if [ ! -f "quicfuzz-fuzz.zip" ]; then
    $HOME/go/bin/go-fuzz-build -tags quicfuzz
fi
exec $HOME/go/bin/go-fuzz -func "Fuzz$1" -workdir "$1"
