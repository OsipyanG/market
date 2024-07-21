package main

import (
	"github.com/OsipyanG/market/services/auth-msv/config"
	"github.com/OsipyanG/market/services/auth-msv/internal/app"
)

func main() {
	conf := config.MustLoad()

	app.Run(conf)
}
