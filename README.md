# GRPC Chat
A simple CLI Client and server communicating through gRPC

## Usage
### Client
```bash
go get github.com/SamuelBagattin/grpc-chat/cmd/grpc-chat
grpc-chat [server_address] [username]
```


### Server
```bash
wget -O grpc-chat-server $(curl https://api.github.com/repos/SamuelBagattin/grpc-chat/releases/latest | jq -r '.assets[1].url') && chmod +x grpc-chat-server
./grpc-chat-server [listening_port]`
```

### Generate files
```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/greeter.proto
```
