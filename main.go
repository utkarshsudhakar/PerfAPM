package main

import (
	"net/http"

	"github.com/utkarshsudhakar/PerfAPM/apmservice"
)

func main() {

	apmservice.Init()
	http.ListenAndServe(":4047", nil)
}
