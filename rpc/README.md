**pull git submodules and enter this dir first**

#### app.proto

```shell
protoc --go_out=./appProto --go-grpc_out=./appProto ./protos/app.proto
```