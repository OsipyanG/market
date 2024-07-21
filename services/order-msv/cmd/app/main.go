package main

import (
	"github.com/OsipyanG/market/services/order-msv/config"
	"github.com/OsipyanG/market/services/order-msv/internal/app"
)

func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}
