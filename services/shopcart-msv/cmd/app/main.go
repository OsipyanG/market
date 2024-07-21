package main

import (
	"github.com/OsipyanG/market/services/shopcart-msv/config"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/app"
)

func main() {
	conf := config.MustLoad()

	app.Run(conf)
}
