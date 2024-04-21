package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Only Post requests are allowed", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")

	checker(err)

	defer file.Close()

	log.Printf("Uploaded file: %+v\n", header.Filename)

	dst, err := os.Create("./uploads/" + header.Filename)
	checker(err)

	defer dst.Close()
	_, err = io.Copy(dst, file)
	checker(err)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File Successfully uploaded"))

}

func checker(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	var port string = "80"
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", uploadHandler)
	log.Println("Starting the server on :" + port)
	err := http.ListenAndServe(":"+port, nil)
	checker(err)
}
