all: make-dir gen-pb build-srv build-cli
	echo "Success!"

make-dir:
	mkdir -p ./proto/api/pb
	mkdir -p ./out

gen-pb:
	protoc --go_out=./proto/api/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./proto/api/pb \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./proto \
		api.proto

build-srv:
	go build -o ./out/server ./server/main.go

build-cli:
	go build -o ./out/client ./client/main.go
