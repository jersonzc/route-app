## gRPC
Generating client and server code.
```sh
protoc --go_out=. --go_opt=module=route-app \
       --go-grpc_out=. --go-grpc_opt=module=route-app \
       proto/route.proto
```

## Build
```
go build -o server.exe cmd/server/*
go build -o client.exe cmd/client/*
```
