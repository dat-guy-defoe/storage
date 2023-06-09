package mongo

import (
	"fmt"
	mongodb "github.com/dat-guy-defoe/storage/internal/repo/mongo"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	defer func(start time.Time) {
		fmt.Printf("uploadFile exec in: %v\n", time.Since(start))
	}(time.Now())

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error retrieving the file")
		log.Println(err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)

		return
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error saving the file", http.StatusInternalServerError)

		return
	}

	db, err := mongodb.GetClient()
	if err != nil {
		log.Println(err)
		http.Error(w, "Repository access denied", http.StatusInternalServerError)

		return
	}

	err = db.PutFile(handler.Filename, fileBytes)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	defer func(start time.Time) {
		fmt.Printf("downloadFile exec in: %v\n", time.Since(start))
	}(time.Now())

	vars := mux.Vars(r)
	filename := vars["fileName"]

	db, err := mongodb.GetClient()
	if err != nil {
		log.Println(err)
		http.Error(w, "Repository access denied", http.StatusInternalServerError)

		return
	}

	data, err := db.GetFile(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get the file", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Failed to get the file", http.StatusInternalServerError)
	}
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	defer func(start time.Time) {
		fmt.Printf("deleteFile exec in: %v\n", time.Since(start))
	}(time.Now())

	vars := mux.Vars(r)
	filename := vars["fileName"]

	db, err := mongodb.GetClient()
	if err != nil {
		log.Println(err)
		http.Error(w, "Repository access denied", http.StatusInternalServerError)

		return
	}

	err = db.DeleteFile(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete file %s", filename), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
