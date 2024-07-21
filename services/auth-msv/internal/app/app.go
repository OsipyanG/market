package app

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	authpb "github.com/OsipyanG/market/protos/auth"
	"github.com/OsipyanG/market/services/auth-msv/config"
	"github.com/OsipyanG/market/services/auth-msv/internal/cache/memcached"
	adminconroller "github.com/OsipyanG/market/services/auth-msv/internal/controller/admin"
	authconroller "github.com/OsipyanG/market/services/auth-msv/internal/controller/auth"
	adminservice "github.com/OsipyanG/market/services/auth-msv/internal/service/admin"
	authservice "github.com/OsipyanG/market/services/auth-msv/internal/service/auth"
	"github.com/OsipyanG/market/services/auth-msv/internal/storage/postgres"
	"google.golang.org/grpc"
)

func Run(conf *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	storage, err := postgres.New(ctx, &conf.Storage)
	if err != nil {
		slog.Error("error while connecting to database: ", "err", err)

		return
	}
	defer storage.Close()

	refreshTokensCache, err := memcached.New(&conf.Memcache, &conf.Refresh)
	if err != nil {
		slog.Error("error while connecting to memcached: ", "err", err)

		return
	}
	defer refreshTokensCache.Close()

	authConf := &config.AuthServiceConfig{
		Security: conf.Security,
		Access:   conf.Access,
	}

	authService := authservice.New(authConf, storage, refreshTokensCache)
	authController := authconroller.New(authService)

	adminService := adminservice.New(storage)
	adminController := adminconroller.New(adminService)

	address := net.JoinHostPort(conf.GRPCServer.Host, conf.GRPCServer.Port)

	listen, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("failed to listen: ", "err", err)

		return
	}

	server := grpc.NewServer()

	authpb.RegisterAuthServer(server, authController)
	authpb.RegisterAuthAdminServer(server, adminController)

	go func() {
		slog.Info("[Auth]: start listening at ", "addr", listen.Addr())

		err := server.Serve(listen)
		if err != nil {
			slog.Error("[Auth]: failed to serve: ", "err", err)
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
	slog.Info("the server has shut down.")
}
