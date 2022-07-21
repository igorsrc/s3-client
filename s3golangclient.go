package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"s3golangclient/config"
	"s3golangclient/fs"
	"s3golangclient/handler"
)

func init() {
	config.InitConfig("config.yaml")
}

func main() {
	cfg := config.Get()
	serverAddr := cfg.S3.Host + ":" + cfg.S3.Port

	fs.InitClient(serverAddr, cfg.S3.Access, cfg.S3.Secret, cfg.S3.Region)
	fs.InitBucket(cfg.S3.Bucket)

	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/upload", handler.UploadHandler).Methods("POST")
	appRouter.HandleFunc("/download", handler.DownloadHandler).Methods("GET").Queries("key", "{key}")

	srv := &http.Server{Handler: appRouter, Addr: ":8080"}
	log.Fatal(srv.ListenAndServe())
}
