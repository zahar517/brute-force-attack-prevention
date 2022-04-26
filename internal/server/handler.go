package grpcserver

import (
	context "context"

	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Login(ctx context.Context, in *LoginRequest) (*LoginResponse, error) {
	ok, err := s.app.Login(ctx, in.Login, in.Password, in.Ip)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Ok: ok}, nil
}

func (s *Server) AddToWhitelist(ctx context.Context, in *SubnetRequest) (*empty.Empty, error) {
	if err := s.app.AddToWhitelist(ctx, in.Subnet); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *Server) RemoveFromWhitelist(ctx context.Context, in *SubnetRequest) (*empty.Empty, error) {
	if err := s.app.RemoveFromWhitelist(ctx, in.Subnet); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *Server) AddToBlacklist(ctx context.Context, in *SubnetRequest) (*empty.Empty, error) {
	if err := s.app.AddToBlacklist(ctx, in.Subnet); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *Server) RemoveFromBlacklist(ctx context.Context, in *SubnetRequest) (*empty.Empty, error) {
	if err := s.app.RemoveFromBlacklist(ctx, in.Subnet); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
