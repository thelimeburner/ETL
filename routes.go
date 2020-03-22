package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func fetchAcceptType(header string) string {
	if strings.Contains(header, "application/json") {
		return "json"
	}
	return "html"

}

//BasicAuth handles authentication for username and password
func BasicAuth(handler http.HandlerFunc, perm, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		users := LogStore.fetchUserAuth(perm)

		user, pass, ok := r.BasicAuth()

		for _, u := range users {
			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(u.User)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(u.Password)) != 1 {
				continue
			} else {
				handler(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))

		//MoDIFY
		return
	}
}

//fetchBrowserCounts
func handleBrowserCount(w http.ResponseWriter, r *http.Request) {
	browserData := LogStore.fetchBrowserData()

	browserStats := make(map[string]int)

	for _, v := range browserData {
		if _, ok := browserStats[v.Browser]; !ok {
			browserStats[v.Browser] = 0
		}
		browserStats[v.Browser] += v.Count

	}

	jOut, err := json.Marshal(browserStats)
	if err != nil {
		fmt.Println("Error Unmarshalling data", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jOut))
}

//fetchBrowserCounts
func handleVisitorCount(w http.ResponseWriter, r *http.Request) {

	visitorData := LogStore.fetchVisitorData()

	visitorStats := make(map[string]int)

	for _, v := range visitorData {

		if _, ok := visitorStats[v.Key]; !ok {
			visitorStats[v.Key] = 0
		}
		visitorStats[v.Key] += v.Count

	}

	jOut, err := json.Marshal(visitorStats)
	if err != nil {
		fmt.Println("Error Unmarshalling data", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jOut))
}

//handleServeUploadPage serves the static html file
func handleServeUploadPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/upload.html")
}

//handleUploadLog handles uploading of log file and triggers etl pipeline
func handleUploadLog(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Received Uploaded File: %+v\n", handler.Filename)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("File Contents", string(fileBytes))

	// return that we have successfully uploaded our file!
	processLogFile(fileBytes)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Log File Uploaded Successfully")
}
