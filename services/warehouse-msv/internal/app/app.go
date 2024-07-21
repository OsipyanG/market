package app

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	warehousepb "github.com/OsipyanG/market/protos/warehouse"
	"github.com/OsipyanG/market/services/warehouse-msv/config"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/controller"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/repository/postgres"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pgPool, err := postgres.NewClient(ctx, &cfg.Postgres)
	if err != nil {
		slog.Error("error while connecting to database: ", "err", err)

		return
	}
	defer pgPool.Close()

	warehouseRepository := postgres.NewWarehouseRepository(pgPool)
	warehouseService := service.NewWarehouseService(warehouseRepository)
	warehouseController := controller.NewWarehouseController(warehouseService)

	catalogRepository := postgres.NewCatalogRepository(pgPool)
	catalogService := service.NewCatalogService(catalogRepository)
	catalogController := controller.NewCatalogController(catalogService)

	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		slog.Error("cannot listen:", "err", err)

		return
	}

	server := grpc.NewServer()
	warehousepb.RegisterWarehouseServer(server, warehouseController)
	warehousepb.RegisterCatalogServer(server, catalogController)

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
