package handler

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"s3golangclient/config"
	"s3golangclient/fs"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(config.Get().App.MaxUploadSize << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, headers, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if headers.Filename == "" {
		http.Error(w, "missing filename", http.StatusBadRequest)
		return
	}

	defer file.Close()

	log.Printf("uploading file %s of size %s", headers.Filename, headers.Size)

	for key, value := range r.Header {
		log.Printf("%s=%s", key, value)
	}

	key := time.Now().Format("20060102") + "/" + headers.Filename

	_, err = fs.GetS3Client().PutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket:             aws.String(config.Get().S3.Bucket),
			Key:                aws.String(key),
			ContentDisposition: aws.String(r.Header.Get("Content-Disposition")),
			ContentType:        aws.String(r.Header.Get("Content-Type")),
			Body:               file,
		})

	if err != nil {
		http.Error(w, "error uploading file "+headers.Filename, http.StatusInternalServerError)
		return
	}

	log.Printf("File saved with key: %s", key)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	log.Printf("Downloading file by key %s", key)

	getResponse, err := fs.GetS3Client().GetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(config.Get().S3.Bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		http.Error(w, "error downloading file with key "+key, http.StatusInternalServerError)
		return
	}

	if getResponse.ContentDisposition != nil {
		w.Header().Set("Content-Disposition", *getResponse.ContentDisposition)
	}
	if getResponse.ContentType != nil {
		w.Header().Set("Content-Type", *getResponse.ContentType)
	}
	io.Copy(w, getResponse.Body)
}
