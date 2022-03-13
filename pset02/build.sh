rm -rf bin
mkdir bin

# MINER

# 64 bit
GOOS=windows GOARCH=amd64 go build -o bin/miner-amd64.exe
GOOS=darwin GOARCH=amd64 go build -o bin/miner-amd64-darwin
GOOS=linux GOARCH=amd64 go build -o bin/miner-amd64-linux

# SERVER

# 64 bit
GOOS=windows GOARCH=amd64 go build -o bin/server-amd64.exe server/server.go
GOOS=darwin GOARCH=amd64 go build -o bin/server-amd64-darwin server/server.go
GOOS=linux GOARCH=amd64 go build -o bin/server-amd64-linux server/server.go