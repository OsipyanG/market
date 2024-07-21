package app

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	orderpb "github.com/OsipyanG/market/protos/order"
	shopcartpb "github.com/OsipyanG/market/protos/shopcart"
	warehousepb "github.com/OsipyanG/market/protos/warehouse"
	"github.com/OsipyanG/market/services/order-msv/config"
	"github.com/OsipyanG/market/services/order-msv/internal/controller"
	"github.com/OsipyanG/market/services/order-msv/internal/repository/postgres"
	"github.com/OsipyanG/market/services/order-msv/internal/service"
	shopcart "github.com/OsipyanG/market/services/order-msv/internal/transport/shopcart/grpc"
	warehouse "github.com/OsipyanG/market/services/order-msv/internal/transport/warehouse/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		slog.Error("cannot open database: ", "err", err)

		return
	}
	defer pool.Close()

	shopcartConn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Shopcart.Host, cfg.Shopcart.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("cannot connect to shopcart: ", "err", err)

		return
	}
	defer shopcartConn.Close()

	warehouseConn, err := grpc.NewClient(
		net.JoinHostPort(cfg.Warehouse.Host, cfg.Warehouse.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("cannot connect to warehouse: ", "err", err)

		return
	}
	defer warehouseConn.Close()

	shopcartClient := shopcart.NewClient(shopcartpb.NewUserShopcartClient(shopcartConn))
	warehouseClient := warehouse.NewClient(warehousepb.NewWarehouseClient(warehouseConn))

	orderRepository := postgres.NewOrderRepository(pool)

	orderService := service.NewOrderService(orderRepository, shopcartClient, warehouseClient)
	orderController := controller.NewOrderController(orderService)

	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		slog.Error("cannot listen:", "err", err)

		return
	}

	server := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(server, orderController)

	reflection.Register(server)
	slog.Info("app is listening on ", cfg.GRPC.Host, cfg.GRPC.Port)

	go func() {
		if err := server.Serve(lis); err != nil {
			slog.Error("failed to app: ", "err", err)

			return
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
	slog.Info("server exiting")
}
