package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLogging() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		st, _ := status.FromError(err)
		log.Printf("[UNARY] method=%s duration=%s code=%s error=%v",
			info.FullMethod, time.Since(start), st.Code(), err)
		return resp, err
	}
}

// StreamLogging - Stream 인터셉터: wrappedStream 패턴으로 Send/Recv 로깅
func StreamLogging() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		log.Printf("[STREAM START] method=%s client_stream=%v server_stream=%v",
			info.FullMethod, info.IsClientStream, info.IsServerStream)

		wrapped := &wrappedStream{ServerStream: ss, method: info.FullMethod}
		err := handler(srv, wrapped)

		st, _ := status.FromError(err)
		log.Printf("[STREAM END] method=%s duration=%s code=%s send_count=%d recv_count=%d error=%v",
			info.FullMethod, time.Since(start), st.Code(),
			wrapped.sendCount, wrapped.recvCount, err)
		return err
	}
}

type wrappedStream struct {
	grpc.ServerStream
	method    string
	sendCount int
	recvCount int
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	w.sendCount++
	return w.ServerStream.SendMsg(m)
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		w.recvCount++
	}
	return err
}
