package configs

type GRPCServer struct {
	Addr string `config:"grpc_server_addr"`
}

func newGRPCServer() *GRPCServer {
	return &GRPCServer{
		Addr: ":8081",
	}
}
