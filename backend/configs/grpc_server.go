package configs

type GRPCServer struct {
	Addr string `config:"http_server_addr"`
}

func newGRPCServer() *GRPCServer {
	return &GRPCServer{
		Addr: ":8081",
	}
}
