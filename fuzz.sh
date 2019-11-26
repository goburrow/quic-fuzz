#!/bin/sh
exec $HOME/go/bin/go-fuzz -func "Fuzz$1" -workdir "$1"
