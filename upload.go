package main

import (
	"crypto/md5"
	"fmt"
	"io"
	// "io/ioutil"
	// "log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

const (
	upload_path string = "./uploads"
	server_port string = "2640"
)

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "A Go Web Server")
	io.WriteString(w, "Welcome!")
}

func postCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		printRequest(w, r, true)
	}
}

//上传
func uploadHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		printRequest(w, r, true)
		io.WriteString(w,
			"<html>"+
				"<form action='' method=\"post\" enctype=\"multipart/form-data\">"+
				"<input type=\"file\" name='photos'/><input type=\"submit\" value=\"Upload\"/>"+
				"</form>"+
				"</html>")
	} else if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.Header().Set("Access-Control-Max-Age", "30")
	} else if r.Method == "POST" {
		// ip := strings.Split(r.RemoteAddr, ":")[0]
		printRequest(w, r, false)
		msg, code, err := SaveFileFromRequest(w, r, upload_path)
		if err != nil {
			if code <= 0 {
				code = http.StatusInternalServerError
			}
			Error(w, msg, code, err)
			return
		}
		temp_path := upload_path + "/" + msg
		//md5
		id, _ := FileHashMD5(temp_path)
		real_path := upload_path + "/" + id
		//check exists
		if PathExist(real_path) {
			Error(w, id, http.StatusCreated, nil)
			return
		}
		//response
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.WriteString(w, id)
		fmt.Println("   upload success " + id)
		//rename
		err = os.Rename(temp_path, real_path)
		if err != nil {
			fmt.Println("ERROR TO RENAME: ", err)
		}
	}
}

func main() {

	fmt.Println("starting server!")
	fmt.Println("http://127.0.0.1:" + server_port)

	// http.Post("/media", uploadHandle)
	http.HandleFunc("/", defaultHandle)
	http.HandleFunc("/media", uploadHandle)
	http.HandleFunc("/post", postCreate)

	server := &http.Server{
		Addr: ":" + server_port,
		// Handler:        handler,
		ReadTimeout:    50 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 10,
	}
	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
		return
	}
}

/*--------------UTILS-------------*/

func SaveFileFromRequest(w http.ResponseWriter, r *http.Request, parent string) (string, int, error) {

	fmt.Println("Reading")
	//get file
	file, head, err := r.FormFile("photos")
	if err != nil {
		return "Fail to read file from form", http.StatusInternalServerError, err
	}
	defer file.Close()

	//temp file name
	id := fmt.Sprintf("%x", md5.Sum([]byte(head.Filename)))
	id = id + "-" + head.Filename
	temp_path := parent + "/" + id

	fmt.Println("Creating")
	//create file
	fW, err := os.Create(temp_path)
	if err != nil {
		return "Fail to create file!", http.StatusInternalServerError, err
	}
	defer fW.Close()

	fmt.Println("Coping")
	//save file
	_, err = io.Copy(fW, file)
	if err != nil {
		return "Fail to save file!", http.StatusInternalServerError, err
	}
	return id, http.StatusOK, nil
}

func printRequest(w http.ResponseWriter, r *http.Request, body bool) {
	fmt.Println()
	fmt.Println("------printRequest------")
	fmt.Println("requester:    " + strings.Split(r.RemoteAddr, ":")[0])
	debug(httputil.DumpRequest(r, body))
	fmt.Println()
	fmt.Println()
	fmt.Println("----------END-----------")
	fmt.Println()
}

func Error(w http.ResponseWriter, msg string, code int, err error) {
	if err != nil {
		fmt.Println(err)
	}
	if msg != "" {
		fmt.Println(msg)
	}
	http.Error(w, msg, code)
}

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// func FileHashMD5(file *os.File) (string, error) {
func FileHashMD5(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}
	h := md5.New()
	io.Copy(h, file)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s", data)
	} else {
		fmt.Printf("%s", err)
	}
}
