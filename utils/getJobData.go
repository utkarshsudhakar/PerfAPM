package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/utkarshsudhakar/demo_app/config"
)

func GetJobData(resp *http.Response, client *http.Client, JobURL string, ResourceName string, ldmHeader config.LDMHeader) (config.JobResponse, *http.Client, bool) {

	//fmt.Println(jsonData.InfaToken)
	///URL = "http://irlcmg08.informatica.com:8085/ldmadmin/ldm.resources.core/execute"
	var jobResponse config.JobResponse
	resourceFound := true

	req, _ := http.NewRequest("POST", JobURL, nil)

	q := req.URL.Query()
	q.Add("name", ResourceName)
	q.Add("status", "true")
	q.Add("infaToken", ldmHeader.InfaToken)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/json")
	client.Jar.SetCookies(req.URL, resp.Cookies())
	req.Header.Add("JSESSIONID", ldmHeader.JsessionID)
	//fmt.Println(ldmHeader.JsessionID)
	resp, err := client.Do(req)
	//fmt.Println("in getJob Data")

	if err != nil {
		fmt.Println("inside err")
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &jobResponse)

	//fmt.Println(string(body))
	//fmt.Println(jobResponse)
	if len(body) == 2 {
		resourceFound = false
	}
	return jobResponse, client, resourceFound

}
