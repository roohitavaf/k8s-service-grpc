install proto:

```
brew install protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

install ko:

```
go install github.com/google/ko@latest
```

Add go to path
```
export PATH="$PATH:$(go env GOPATH)/bin"
```




