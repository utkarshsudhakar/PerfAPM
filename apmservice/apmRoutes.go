package apmservice

import "net/http"

func Init() {

	http.HandleFunc("/compareBuild", compareBuild)
	http.HandleFunc("/test", test)
	http.HandleFunc("/createJson", createJson)

}
