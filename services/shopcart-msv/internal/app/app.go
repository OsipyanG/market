package app

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"github.com/OsipyanG/market/protos/shopcart"
	"github.com/OsipyanG/market/services/shopcart-msv/config"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/controller"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/service"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/storage/postgres"
	"google.golang.org/grpc"
)

func Run(config *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	storage, err := postgres.New(ctx, &config.Storage)
	if err != nil {
		slog.Error("error while connecting to database: ", "err", err)

		return
	}
	defer storage.Close()

	shopCartService := service.New(storage)
	shopCartController := controller.New(shopCartService)

	address := net.JoinHostPort(config.GRPCServer.Host, config.GRPCServer.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("failed to listen: ", "err", err)

		return
	}

	server := grpc.NewServer()
	shopcart.RegisterUserShopcartServer(server, shopCartController)

	go func() {
		slog.Info("start listening at ", "addr", lis.Addr())

		err := server.Serve(lis)
		if err != nil {
			slog.Error("failed to serve: ", "err", err)
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
	slog.Info("the server has shut down.")
}
