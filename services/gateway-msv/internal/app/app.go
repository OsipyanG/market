package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	auth "github.com/OsipyanG/market/protos/auth"
	order "github.com/OsipyanG/market/protos/order"
	shopcart "github.com/OsipyanG/market/protos/shopcart"
	catalog "github.com/OsipyanG/market/protos/warehouse"
	"github.com/OsipyanG/market/services/gateway-msv/config"
	"github.com/OsipyanG/market/services/gateway-msv/internal/controller/handler"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	router := gin.Default()

	shopcartConn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Shopcart.Host, cfg.Shopcart.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("connection to shopcart service failed", "error", err)

		return
	}
	defer shopcartConn.Close()

	orderConn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Order.Host, cfg.Order.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("connection to order service failed", "error", err)

		return
	}
	defer orderConn.Close()

	authConn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Auth.Host, cfg.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("connection to auth service failed", "error", err)

		return
	}
	defer authConn.Close()

	catalogConn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Warehouse.Host, cfg.Warehouse.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("connection to catalog service failed", "error", err)

		return
	}
	defer catalogConn.Close()

	authClient := auth.NewAuthClient(authConn)
	authAdminClient := auth.NewAuthAdminClient(authConn)
	authHandler := handler.NewAuthHandler(authClient, authAdminClient)

	catalogClient := catalog.NewCatalogClient(catalogConn)
	catalogHandler := handler.NewCatalogHandler(catalogClient)

	shopcartClient := shopcart.NewUserShopcartClient(shopcartConn)
	shopcartHandler := handler.NewShopcartHandler(shopcartClient)

	orderClient := order.NewOrderServiceClient(orderConn)
	orderHandler := handler.NewOrderHandler(orderClient)

	handler.SetupRouter(router, shopcartHandler, orderHandler, catalogHandler, authHandler)

	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      router.Handler(),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	slog.Info("server exiting")
}
