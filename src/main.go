package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/hamal/src/config"
	"github.com/Dataman-Cloud/hamal/src/router"
	"github.com/Dataman-Cloud/hamal/src/router/middleware"
	_ "github.com/Dataman-Cloud/hamal/src/utils"

	log "github.com/Sirupsen/logrus"
)

func main() {
	configFile := flag.String("config", "config_file", "config file path")
	flag.Parse()
	config.InitConfig(*configFile)

	log.Infof("http listen %s starting...", config.GetConfig().Addr)
	server := &http.Server{
		Addr:           config.GetConfig().Addr,
		Handler:        router.Router(middleware.Authenticate),
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("http listen server error: %v", err)
	}

}
