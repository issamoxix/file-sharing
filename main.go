package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
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

	http.Redirect(w, r, "/", http.StatusFound)
	w.Write([]byte("File Successfully uploaded"))

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
