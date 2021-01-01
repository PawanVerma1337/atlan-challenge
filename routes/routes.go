package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pawanverma1337/atlan-challenge/db"
	"github.com/pawanverma1337/atlan-challenge/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// intializeRoutes func initializes all the routes.
func initializeRoutes(r *mux.Router) {
	r.HandleFunc("/upload", postFileUpload).Methods("POST")               // Post file handler
	r.HandleFunc("/upload/terminate", terminateFileUpload).Methods("GET") // Terminate
	r.HandleFunc("/upload/resume", resumeFileUpload).Methods("PATCH")     // Resume
	r.HandleFunc("/upload/status", statusFileUpload).Methods("GET")       // Status
}

func postFileUpload(w http.ResponseWriter, req *http.Request) {
	ul, err := strconv.Atoi(req.Header.Get("Upload-Length"))
	if err != nil {
		e := "Improper upload length"
		log.Printf("%s %s", e, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e))
		return
	}
	log.Printf("upload length %d\n", ul)

	io := 0
	uc := false

	f := db.File{
		ID:                 primitive.NewObjectID(),
		FileUploadOffset:   io,
		FileUploadLength:   ul,
		FileUploadComplete: uc,
	}

	err = f.CreateFile()
	if err != nil {
		e := "Error creating file in DB"
		log.Printf("%s %s\n", e, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filePath := path.Join("files", f.ID.Hex())
	file, err := os.Create(filePath)
	if err != nil {
		e := "Error creating file in filesystem"
		log.Printf("%s %s\n", e, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	w.Header().Set("Location", fmt.Sprintf("localhost:8080/files/%s", f.ID.Hex()))
	w.WriteHeader(http.StatusCreated)
	return
}

func terminateFileUpload(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	id := query["id"][0]
	idO, err := primitive.ObjectIDFromHex(id)
	util.CheckError(err)
	f := db.File{
		ID: idO,
	}
	err = f.UpdateFile(bson.D{{"$set", bson.D{{"file_terminate", true}}}})
	util.CheckError(err)

	w.WriteHeader(http.StatusAccepted)
}

func resumeFileUpload(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	id := query["id"][0]
	idO, err := primitive.ObjectIDFromHex(id)
	util.CheckError(err)
	file := db.File{
		ID: idO,
	}
	err = file.GetFile()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if file.FileUploadComplete == true {
		e := "Upload already completed"
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(e))
		return
	}
	off, err := strconv.Atoi(req.Header.Get("Upload-Offset"))
	if err != nil {
		log.Println("Improper upload offset", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Upload offset %d\n", off)
	if file.FileUploadOffset != off {
		e := fmt.Sprintf("Expected Offset %d got offset %d", file.FileUploadOffset, off)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(e))
		return
	}

	log.Println("Content length is", req.Header.Get("Content-Length"))
	clh := req.Header.Get("Content-Length")
	cl, err := strconv.Atoi(clh)
	if err != nil {
		log.Println("unknown content length")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cl != (file.FileUploadLength - file.FileUploadOffset) {
		e := fmt.Sprintf("Content length doesn't not match upload length.Expected content length %d got %d", file.FileUploadLength-file.FileUploadOffset, cl)
		log.Println(e)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e))
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Received file partially %s\n", err)
		log.Println("Size of received file ", len(body))
	}
	fp := path.Join("files", file.ID.Hex())
	fmt.Println(fp)
	f, err := os.OpenFile(fp, os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("unable to open file %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	n, err := f.WriteAt(body, int64(off))
	if err != nil {
		log.Printf("unable to write %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("number of bytes written ", n)
	no := file.FileUploadOffset + n
	file.FileUploadOffset = no

	uo := strconv.Itoa(file.FileUploadOffset)
	w.Header().Set("Upload-Offset", uo)
	if file.FileUploadOffset == file.FileUploadLength {
		log.Println("upload completed successfully")
		file.FileUploadComplete = true
	}

	err = file.UpdateFile(bson.D{{"$set", bson.D{{"file_completed", false}, {"file_offset", file.FileUploadOffset}}}})
	if err != nil {
		log.Println("Error while updating file", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

	return
}

func statusFileUpload(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	id := query["id"][0]
	idO, err := primitive.ObjectIDFromHex(id)
	util.CheckError(err)
	f := db.File{
		ID: idO,
	}
	err = f.GetFile()
	util.CheckError(err)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

// Run fuction to run the initialization of routes
// from the main.go file.
func Run(r *mux.Router) {
	initializeRoutes(r)
}
