package main

import (
	"net/http"

	"./apmservice"
)

func main() {

	apmservice.Init()
	http.ListenAndServe(":4047", nil)
}
