package main

import (
	"fmt"
	"regexp"

	"github.com/anaskhan96/soup"
)

func main() {

	http.HandleFunc("/getVin", getVin)
	http.ListenAndServe(":4047", nil)

	resp, err := soup.Get("https://carstory.com/cars/make-dodge/model-caravan")
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println(resp)
	r, _ := regexp.Compile("vin\":\"(.*?)\"")
	f := r.FindAllStringSubmatch(resp, -1)
	for i := range f {
		fmt.Println(f[i][1])

	}

}
