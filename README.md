## gRPC
Generating client and server code.
```sh
protoc --go_out=. --go_opt=module=route-app \
       --go-grpc_out=. --go-grpc_opt=module=route-app \
       proto/route.proto
```
