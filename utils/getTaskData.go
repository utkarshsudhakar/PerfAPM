package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/utkarshsudhakar/demo_app/config"
)

func GetTaskData(InfaToken string, jobID string, taskID string, TaskURL string, resp *http.Response, client *http.Client) config.TaskResponse {

	req, _ := http.NewRequest("GET", TaskURL, nil)
	var taskResponse config.TaskResponse

	q := req.URL.Query()
	q.Add("taskid", taskID)
	q.Add("jobid", jobID)
	q.Add("infaToken", InfaToken)
	req.URL.RawQuery = q.Encode()

	client.Jar.SetCookies(req.URL, resp.Cookies())

	resp, err := client.Do(req)

	if err != nil {

		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &taskResponse)
	//fmt.Println(taskResponse)

	return taskResponse

}
