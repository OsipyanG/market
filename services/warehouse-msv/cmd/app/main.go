package main

import (
	"github.com/OsipyanG/market/services/warehouse-msv/config"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/app"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}
