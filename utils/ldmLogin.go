package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func LDMLogin(URL string, client *http.Client, User string, Pass string) (*http.Response, *http.Client) {

	values := map[string]string{"user": User, "password": Pass, "namespace": "Native"}

	jsonValue, _ := json.Marshal(values)
	//pass the values to the request's body

	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {

		fmt.Println(resp)
	}

	return resp, client
}
