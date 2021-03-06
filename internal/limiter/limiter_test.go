package limiter_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zahar517/brute-force-attack-prevention/internal/limiter"
)

const (
	interval = 5
	limit    = 10
)

type LimiterTestSuite struct {
	suite.Suite
	l *limiter.LeakyBucket
}

func (s *LimiterTestSuite) SetupTest() {
	s.l = limiter.New(limit, limit+1, limit+2, interval)
}

func TestLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}

func (s *LimiterTestSuite) TestAddTrue() {
	var i int64
	for ; i < limit; i++ {
		res := s.l.Add("login", "pass", "ip")
		s.Require().True(res)
	}
}

func (s *LimiterTestSuite) TestAddFalse() {
	var i int64
	for ; i < limit; i++ {
		s.l.Add("login", "pass", "ip")
	}

	res := s.l.Add("login", "pass", "ip")
	s.Require().False(res)
}

func (s *LimiterTestSuite) TestReset() {
	var i int64
	for ; i < limit; i++ {
		s.l.Add("login", "pass", "ip")
	}

	s.l.Reset("login", "ip")
	res := s.l.Add("login", "pass", "ip")
	s.Require().True(res)
}
