package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/noydev/ggmapapi-test/domain"
	"github.com/noydev/ggmapapi-test/services"
	util "github.com/noydev/ggmapapi-test/utils"
	"github.com/noydev/ggmapapi-test/utils/logger"
)

func sanityCheck() {
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Fatal(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}

func Start() {

	// sanityCheck()
	cfg, err := util.LoadConfig()
	if err != nil {
		os.Exit(1)
	}
	domain := domain.InitDomains(cfg)

	router := mux.NewRouter()

	mh := MapHandler{services.NewMapService(domain)}

	// define routes
	router.
		HandleFunc("/map/restaurant", mh.GetRestaurantByLocation).
		Methods(http.MethodGet).
		Name("GetRestaurantByLocation")

	// starting server
	address := cfg.Cfg.GetString("app.address")
	port := cfg.Cfg.GetString("app.port")
	fmt.Println(address, port)
	logger.Info(fmt.Sprintf("Starting server on %s:%s ...", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))

}
