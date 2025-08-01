## Route App
A simple route mapping application that lets clients get information about
features on their route, create a summary of their route, and exchange
route information such as traffic updates with the server and other clients.

### Setup
Generate client and server code.
```sh
protoc --go_out=. --go_opt=module=route-app \
       --go-grpc_out=. --go-grpc_opt=module=route-app \
       proto/route.proto
```

### Build
```
go build -o server.exe cmd/server/*
go build -o client.exe cmd/client/*
```
