package main

import (
	"log"

	"strconv"

	"github.com/elBroom/meteo/app/config"
	"github.com/elBroom/meteo/app/db"
	"github.com/elBroom/meteo/app/router"
	"github.com/valyala/fasthttp"
)

func main() {
	sql_connect := db.Sql_connect()
	defer sql_connect.Close()

	cfg := config.GetApp()
	go cfg.Hub.Run()

	router := router.Routing()
	log.Printf("Start server: 127.0.0.1:%d\n", cfg.Port)
	log.Fatal(fasthttp.ListenAndServe(":"+strconv.Itoa(cfg.Port),
		fasthttp.TimeoutHandler(router.Handler, config.RequestWaitInQueueTimeout, "timeout")))
}
