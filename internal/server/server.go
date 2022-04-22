//go:generate protoc -I ../../api/ BFAPTool.proto --go_out=. --go-grpc_out=.

package grpcserver

import (
	context "context"
	"net"

	grpc "google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	logger Logger
	app    Application
	addr   string
	UnimplementedBFAPToolServer
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface {
	Login(ctx context.Context, login, password, ip string) error
	ResetBuket(ctx context.Context, login, ip string) error
	AddToBlacklist(ctx context.Context, subnet string) error
	RemoveFromBlacklist(ctx context.Context, subnet string) error
	AddToWhitelist(ctx context.Context, subnet string) error
	RemoveFromWhitelist(ctx context.Context, subnet string) error
}

func NewServer(logger Logger, app Application, host string, port string) *Server {
	addr := net.JoinHostPort(host, port)
	return &Server{nil, logger, app, addr, UnimplementedBFAPToolServer{}}
}

func (s *Server) Start(_ context.Context) error {
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(loggerInterceptor(s.logger)))
	RegisterBFAPToolServer(server, s)

	s.server = server

	return server.Serve(lsn)
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}
