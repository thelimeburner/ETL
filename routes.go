package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
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
func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		users := LogStore.fetchUserAuth("read")

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
		return
	}
}

//fetchBrowserCounts
func handleBrowserCount(w http.ResponseWriter, r *http.Request) {
	browserData := LogStore.fetchBrowserData()

	browserStats := make(map[string]int)
	//fmt.Println(browserStats)
	for _, v := range browserData {
		if _, ok := browserStats[v.Browser]; !ok {
			browserStats[v.Browser] = 0
		}
		browserStats[v.Browser] += v.Count

	}
	// for k, v := range browserStats {
	// 	fmt.Println(k, ": ", v)

	// }
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
