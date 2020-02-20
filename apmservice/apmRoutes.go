package apmservice

import "net/http"

func Init() {

	http.HandleFunc("/compareBuild", compareBuild)
	http.HandleFunc("/compareRelease", compareRelease)
	http.HandleFunc("/test", test)
	http.HandleFunc("/createJson", createJson)
	http.HandleFunc("/htmlReport/", htmlReport)

}
