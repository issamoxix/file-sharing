package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving file from form data:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadedFile, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		fmt.Println("Error creating file on server:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(uploadedFile, file)
	if err != nil {
		fmt.Println("Error copying file data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("File uploaded successfully:", handler.Filename)
	w.WriteHeader(http.StatusOK)

}

func GetLocalIPs() []net.IP {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	checker(err)

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips
}

func checker(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func WriteAddress(address string) {
	os.WriteFile("./static/address.txt", []byte(address), 0660)
}

func main() {
	var port string = "80"
	localIps := GetLocalIPs()
	localip := localIps[0].String()
	WriteAddress(localip)

	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err := os.Mkdir("./uploads", 0755)
		checker(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", uploadHandler)
	log.Println("Starting the server on :" + port)
	exec.Command("cmd", []string{"/c", "start", "http://" + localip}...).Start()
	err := http.ListenAndServe(":"+port, nil)
	checker(err)
}
