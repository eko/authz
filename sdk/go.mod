module github.com/eko/authz/sdk

go 1.20

require (
	github.com/eko/authz/backend v0.8.4
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29
	google.golang.org/grpc v1.54.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/eko/authz/backend => ../backend
