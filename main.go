package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/girishg4t/app_invite_service/pkg/db"
	"github.com/girishg4t/app_invite_service/pkg/logging"
	"github.com/girishg4t/app_invite_service/pkg/middleware"
	"github.com/girishg4t/app_invite_service/pkg/model"
	"github.com/girishg4t/app_invite_service/pkg/repo"
	"github.com/girishg4t/app_invite_service/pkg/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var tempDir string = "./tmp"

func main() {
	logger := logging.GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Error("error loading .env file")
	}

	dbCfg := model.Database{
		DataSource: filepath.Join(tempDir, "pulseid.db"),
		Debug:      len(os.Getenv("DEBUG")) != 0,
		Schema:     "./scripts",
		Type:       "sqlite3",
	}

	conn, err := db.ConnectToDb(dbCfg)
	if err != nil {
		panic(err)
	}
	repoUserProc := repo.NewUserProcessor(conn, logging.GetLogger())
	s := service.NewUserService(repoUserProc)

	repoAppTokenProc := repo.NewAppTokenProcessor(conn)
	at := service.NewAppTokenService(repoAppTokenProc)
	r := mux.NewRouter()
	publicRoute := r.PathPrefix("/public/api/v1").Subrouter()
	publicRoute.Path("/login").HandlerFunc(s.Login).Methods(http.MethodPost)
	publicRoute.Path("/validatetoken/{appToken}").HandlerFunc(at.ValidateToken).Methods(http.MethodGet)
	// rate limiter is used from this link https://stackoverflow.com/questions/60406965/rate-limit-http-requests-based-on-host-address
	publicRoute.Use(middleware.Limit)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.Authenticate)
	api.Path("/genToken").HandlerFunc(at.GenToken).Methods(http.MethodGet)
	api.Path("/getAllToken").HandlerFunc(at.GetAllAppToken).Methods(http.MethodGet)
	api.Path("/invalidateToken").HandlerFunc(at.InvalidateToken).Methods(http.MethodPatch)
	port := os.Getenv("PORT")
	log.Println("Running local on port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r))
}
