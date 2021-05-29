Fuzzing for https://github.com/goburrow/quic

```
# Optionally disable go module to test against local development
export GO111MODULE=off

go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build

~/go/bin/go-fuzz-build -tags quicfuzz
~/go/bin/go-fuzz -func "FuzzClientInitial" -workdir ClientInitial
```
