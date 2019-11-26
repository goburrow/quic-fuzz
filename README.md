Fuzzing for https://github.com/goburrow/quic

```
go-fuzz-build
go-fuzz -func "FuzzClientInitial" -workdir ClientInitial
```