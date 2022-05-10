// +build integration

package integration_test

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"github.com/zahar517/brute-force-attack-prevention/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BFGrpcSuite struct {
	suite.Suite
	ctx        context.Context
	conn       *grpc.ClientConn
	client     grpcserver.BFAPToolClient
	loginLimit int64
	passLimit  int64
	ipLimit    int64
}

func (s *BFGrpcSuite) SetupSuite() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	grpcHost := os.Getenv("GRPC_HOST")
	grpcPort := os.Getenv("GRPC_PORT")
	loginLimit, err := strconv.ParseInt(os.Getenv("LOGIN_LIMIT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	passLimit, err := strconv.ParseInt(os.Getenv("PASSWORD_LIMIT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	ipLimit, err := strconv.ParseInt(os.Getenv("IP_LIMIT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(
		net.JoinHostPort(grpcHost, grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.conn = conn
	s.client = grpcserver.NewBFAPToolClient(conn)
	s.loginLimit = loginLimit
	s.passLimit = passLimit
	s.ipLimit = ipLimit
}

func (s *BFGrpcSuite) TearDownSuite() {
	s.conn.Close()
}

func TestBFGrpcSuite(t *testing.T) {
	suite.Run(t, new(BFGrpcSuite))
}

func (s *BFGrpcSuite) TestLoginTrue() {
	req := &grpcserver.LoginRequest{Login: "abc", Password: "def", Ip: "192.168.1.3"}

	var i int64
	for ; i < s.loginLimit; i++ {
		res, err := s.client.Login(s.ctx, req)
		s.Require().NoError(err)
		s.Require().True(res.Ok)
	}
}

func (s *BFGrpcSuite) TestLoginFalse() {
	req := &grpcserver.LoginRequest{Login: "cde", Password: "efg", Ip: "192.168.1.5"}

	var i int64
	for ; i < s.loginLimit; i++ {
		res, err := s.client.Login(s.ctx, req)
		s.Require().NoError(err)
		s.Require().True(res.Ok)
	}

	res, err := s.client.Login(s.ctx, req)
	s.Require().NoError(err)
	s.Require().False(res.Ok)
}

func (s *BFGrpcSuite) TestReset() {
	login := "efg"
	password := "ghf"
	ip := "192.168.1.6"
	req := &grpcserver.LoginRequest{Login: login, Password: password, Ip: ip}

	var i int64
	for ; i < s.loginLimit; i++ {
		res, err := s.client.Login(s.ctx, req)
		s.Require().NoError(err)
		s.Require().True(res.Ok)
	}

	_, err := s.client.ResetBuket(s.ctx, &grpcserver.ResetBucketRequest{Login: login, Ip: ip})
	s.Require().NoError(err)

	res, err := s.client.Login(s.ctx, req)
	s.Require().NoError(err)
	s.Require().True(res.Ok)
}

func (s *BFGrpcSuite) TestIntervalReset() {
	req := &grpcserver.LoginRequest{Login: "ikl", Password: "lmn", Ip: "192.168.1.7"}

	var i int64
	for ; i < s.loginLimit; i++ {
		res, err := s.client.Login(s.ctx, req)
		s.Require().NoError(err)
		s.Require().True(res.Ok)
	}

	time.Sleep(time.Second * 60)

	res, err := s.client.Login(s.ctx, req)
	s.Require().NoError(err)
	s.Require().True(res.Ok)
}
