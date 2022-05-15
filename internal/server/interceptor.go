package grpcserver

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func loggerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		var addr string
		var userAgent string

		start := time.Now()

		if p, ok := peer.FromContext(ctx); ok {
			addr = p.Addr.String()
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			userAgents := md.Get("user-agent")
			if len(userAgents) > 0 {
				userAgent = userAgents[0]
			}
		}

		h, err := handler(ctx, req)

		st, _ := status.FromError(err)

		str := fmt.Sprintf(
			"%v [%v] %v %v %v %v",
			addr,
			start.String(),
			info.FullMethod,
			st.Code(),
			time.Since(start),
			userAgent,
		)

		if st.Code() == codes.OK {
			logger.Info(str)
		} else {
			logger.Error(str)
		}

		return h, err
	}
}
