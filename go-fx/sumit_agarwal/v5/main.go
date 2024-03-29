package main

import (
	"context"
	"net"
	"net/http"

	httpServer "github.com/kenshin579/tutorials-go/go-fx/sumit_agarwal/v5/http"
	"github.com/kenshin579/tutorials-go/go-fx/sumit_agarwal/v5/loggerfx"
	pb "github.com/kenshin579/tutorials-go/go-fx/sumit_agarwal/v5/proto"
	rpcServer "github.com/kenshin579/tutorials-go/go-fx/sumit_agarwal/v5/rpc"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	fx.New(
		fx.Provide(http.NewServeMux),
		fx.Provide(rpcServer.New),
		fx.Invoke(httpServer.New),
		fx.Invoke(registerHooks),
		loggerfx.Module,
	).Run()
}

func registerHooks(
	lifecycle fx.Lifecycle, mux *http.ServeMux, logger *zap.SugaredLogger, rpcServer rpcServer.Handler,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {

				// rpc server
				lis, err := net.Listen("tcp", ":8081")
				if err != nil {
					logger.Fatalf("failed to listen: %v", err)
				}
				var opts []grpc.ServerOption
				grpcServer := grpc.NewServer(opts...)
				pb.RegisterUsersServer(grpcServer, rpcServer)
				go grpcServer.Serve(lis)

				// start the http server
				logger.Info("Listening on localhost:8080 for HTTP requests")
				go http.ListenAndServe(":8080", mux) // we will look into how to gracefully handle these errors later

				return nil
			},
			OnStop: func(context.Context) error {
				return logger.Sync()
			},
		},
	)
}
