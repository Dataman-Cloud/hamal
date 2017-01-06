package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/Dataman-Cloud/hamal/src/config"

	log "github.com/Sirupsen/logrus"
)

type HamalService struct {
	Projects map[string]string
	Client   *http.Client
}

func InitHamalService() *HamalService {
	u, err := url.Parse(config.GetConfig().SwanAddr)
	if err != nil {
		log.Fatalf("invalid swan url: %s", config.GetConfig().SwanAddr)
		return nil
	}
	return &HamalService{
		Projects: make(map[string]string),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
