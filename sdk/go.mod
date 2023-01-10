module github.com/eko/authz/sdk

go 1.19

require (
	github.com/eko/authz/backend v0.0.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	golang.org/x/exp v0.0.0-20221208152030-732eee02a75a
	google.golang.org/grpc v1.51.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v0.9.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20210226172003-ab064af71705 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/eko/authz/backend => ../backend
