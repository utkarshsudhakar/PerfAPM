package apmservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/sjson"
	"github.com/utkarshsudhakar/demo_app/config"
	"github.com/utkarshsudhakar/demo_app/utils"
)

func test(w http.ResponseWriter, r *http.Request) {

	body := config.Body{ResponseCode: 200, Message: "OK"}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

}

func compareRelease(w http.ResponseWriter, r *http.Request) {

	oldBuildNum := r.URL.Query().Get("oldBuildNum") //427.4
	oldRelease := r.URL.Query().Get("oldRelease")   // "10.2.2"
	newBuildNum := r.URL.Query().Get("newBuildNum") //427.5
	newRelease := r.URL.Query().Get("newRelease")
	Hostname := r.URL.Query().Get("Hostname")
	completeReport := r.URL.Query().Get("completeReport") // True/False
	cc := r.URL.Query().Get("email")

	oldReleaseData, oldBuildDataTask := utils.GetReleaseData(oldBuildNum, oldRelease, Hostname)
	newReleaseData, newBuildDataTask := utils.GetReleaseData(newBuildNum, newRelease, Hostname)

	if (len(newReleaseData) == 0) || (len(oldReleaseData) == 0) {

		utils.RespondWithJSON("BuildNumber/Release not correct or not enough data ", w, r)

	} else {
		subject := fmt.Sprintf("Release Comparison Report for %s & %s", oldRelease, newRelease)
		p := fmt.Sprintf("<body style='background:white'><h3 style='background:#0790bd;color:#fff;padding:5px;text-align:center;border-radius:5px;'> Release Comparison for %s & %s </h3> <br/> <br/>", oldRelease, newRelease)

		for ResourceName, v := range oldReleaseData {

			if _, ok := newReleaseData[ResourceName]; ok {
				p = p + fmt.Sprintf("<div style='background:yellow;text-align:center'><p><b>Resource Summary : %s </p> </b></div>", ResourceName)

				p = p + fmt.Sprintf("<table style='backgound:#fff;border-collapse: collapse;' border = '1' cellpadding = '6'><tbody><tr><td colspan=5 style='text-align:center;background-color:#444;color:white;'><b>Resource Name : %s </b></td></tr><tr><th>Stage</th><th>Release: %s </th><th>Release: %s</th><th>Time Difference</th><th> %% Time Difference</th></tr> ", ResourceName, oldRelease, newRelease)

				for k := range v {

					svOld := oldReleaseData[ResourceName][k].(string)
					svNew := newReleaseData[ResourceName][k].(string)
					timeOld, _ := time.Parse(config.TimeFormat, svOld)
					timeNew, _ := time.Parse(config.TimeFormat, svNew)
					diff := timeNew.Sub(timeOld)

					if ((timeOld.Second() + (timeOld.Minute() * 60) + (timeOld.Hour() * 3600)) > 30) && ((timeNew.Second() + (timeNew.Minute() * 60) + (timeNew.Hour() * 3600)) > 30) {
						if diff <= 0 {
							percDiff := utils.CalcPerc(float64(diff.Seconds()), timeOld)

							p = p + "<tr style='background:#80CA80'><td>" + k + "</td><td>" + svOld + "</td><td>" + svNew + "</td><td>" + diff.String() + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"

						} else {

							percDiff := utils.CalcPerc(float64(diff.Seconds()), timeOld)
							p = p + "<tr style='background:#ff9e82'><td>" + k + "</td><td>" + svOld + "</td><td>" + svNew + "</td><td>" + diff.String() + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
						}

					}
				}
			}

			p = p + "</tbody></table></body><br/><br/>"
			if completeReport == "True" {

				for TaskName, value := range oldBuildDataTask[ResourceName] {
					//fmt.Printf("TaskName is:%s\n", TaskName)

					//ks := reflect.ValueOf(value).MapKeys()
					//fmt.Println(oldBuildDataTask[ResourceName][TaskName])
					val, _ := value.(map[string]interface{})

					if _, ok := newBuildDataTask[ResourceName][TaskName]; ok {
						p = p + fmt.Sprintf("<table style='backgound:#fff;border-collapse: collapse;' border = '1' cellpadding = '6'><tbody><tr><td colspan=5 style='text-align:center;background-color:#7388f9;color:white;'><b>Task Name : %s </b></td></tr><tr><th>Sub Task</th><th>Release: %s </th><th>Release: %s</th><th>Time Difference</th><th> %% Time Difference</th></tr> ", TaskName, oldRelease, newRelease)

						for key := range val {
							//fmt.Println(oldBuildDataTask[TaskName][key])

							svOldTask := oldBuildDataTask[ResourceName][TaskName].(map[string]interface{})[key].(string)
							svNewTask := newBuildDataTask[ResourceName][TaskName].(map[string]interface{})[key].(string)
							timeOldTask, _ := time.Parse(config.TimeFormat, svOldTask)
							timeNewTask, _ := time.Parse(config.TimeFormat, svNewTask)
							diffTask := timeNewTask.Sub(timeOldTask)

							if ((timeOldTask.Second() + (timeOldTask.Minute() * 60) + (timeOldTask.Hour() * 3600)) > 30) && ((timeNewTask.Second() + (timeNewTask.Minute() * 60) + (timeNewTask.Hour() * 3600)) > 30) {
								if diffTask <= 0 {
									percDiffTask := utils.CalcPerc(float64(diffTask.Seconds()), timeOldTask)

									p = p + "<tr style='background:#a7f392'><td>" + key + "</td><td>" + svOldTask + "</td><td>" + svNewTask + "</td><td>" + diffTask.String() + " </td><td>" + strconv.FormatFloat(percDiffTask, 'f', 2, 64) + " %</td></tr>"

								} else {

									percDiffTask := utils.CalcPerc(float64(diffTask.Seconds()), timeOldTask)
									p = p + "<tr style='background:#f99'><td>" + key + "</td><td>" + svOldTask + "</td><td>" + svNewTask + "</td><td>" + diffTask.String() + " </td><td>" + strconv.FormatFloat(percDiffTask, 'f', 2, 64) + " %</td></tr>"
								}
							}
						}

						p = p + "</tbody></table></body><br/><br/>"
					}

				}
			}
		}
		//fmt.Println(p)
		conf := utils.ReadConfig()
		p = p + "<b>Dashboard URL : </b>" + conf.DashboardURL
		utils.SendMail(p, subject, cc)

		//write to file
		fileName := "../htmlReport/EDCPerformance_report_" + oldRelease + "vs" + newRelease + ".html"
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		_, err = f.WriteString(p)
		if err != nil {
			fmt.Println(err)
		}

		//fmt.Println(p)
		utils.RespondWithJSON("Email Sent Successfully", w, r)
	}

}
func compareBuild(w http.ResponseWriter, r *http.Request) {

	oldBuildNum := r.URL.Query().Get("oldBuildNum") //427.4
	newBuildNum := r.URL.Query().Get("newBuildNum") //427.5
	Release := r.URL.Query().Get("Release")
	Hostname := r.URL.Query().Get("Hostname")             // "10.2.2"
	completeReport := r.URL.Query().Get("completeReport") // True/False
	cc := r.URL.Query().Get("email")

	oldBuildData, oldBuildDataTask := utils.GetBuildData(oldBuildNum, Hostname)
	newBuildData, newBuildDataTask := utils.GetBuildData(newBuildNum, Hostname)

	if (len(newBuildData) == 0) || (len(oldBuildData) == 0) {

		utils.RespondWithJSON("BuildNumber not correct or not enough data ", w, r)

	} else {

		subject := fmt.Sprintf("Release : %s - Build Comparison Report for %s & %s", Release, oldBuildNum, newBuildNum)
		p := fmt.Sprintf("<body style='background:white'><h3 style='background:#0790bd;color:#fff;padding:5px;text-align:center;border-radius:5px;'> Build Comparison for %s & %s </h3> <br/><b>Release: %s </b><br/> <br/>", oldBuildNum, newBuildNum, Release)

		for ResourceName, v := range oldBuildData {

			if _, ok := newBuildData[ResourceName]; ok {
				p = p + fmt.Sprintf("<div style='background:yellow;text-align:center'><p><b>Resource Summary : %s </p> </b></div>", ResourceName)

				p = p + fmt.Sprintf("<table style='backgound:#fff;border-collapse: collapse;' border = '1' cellpadding = '6'><tbody><tr><td colspan=5 style='text-align:center;background-color:#444;color:white;'><b>Resource Name : %s </b></td></tr><tr><th>Stage</th><th>Build# %s </th><th>Build# %s</th><th>Time Difference</th><th> %% Time Difference</th></tr> ", ResourceName, oldBuildNum, newBuildNum)

				for k := range v {
					//fmt.Println(oldBuildData[ResourceName][k])
					svOld := oldBuildData[ResourceName][k].(string)
					svNew := newBuildData[ResourceName][k].(string)
					timeOld, _ := time.Parse(config.TimeFormat, svOld)
					timeNew, _ := time.Parse(config.TimeFormat, svNew)
					diff := timeNew.Sub(timeOld)

					if ((timeOld.Second() + (timeOld.Minute() * 60) + (timeOld.Hour() * 3600)) > 30) && ((timeNew.Second() + (timeNew.Minute() * 60) + (timeNew.Hour() * 3600)) > 30) {
						if diff <= 0 {
							percDiff := utils.CalcPerc(float64(diff.Seconds()), timeOld)

							p = p + "<tr style='background:#a7f392'><td>" + k + "</td><td>" + svOld + "</td><td>" + svNew + "</td><td>" + diff.String() + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"

						} else {

							percDiff := utils.CalcPerc(float64(diff.Seconds()), timeOld)
							p = p + "<tr style='background:#f99'><td>" + k + "</td><td>" + svOld + "</td><td>" + svNew + "</td><td>" + diff.String() + " </td><td>" + strconv.FormatFloat(percDiff, 'f', 2, 64) + " %</td></tr>"
						}
					}

				}
			}
			p = p + "</tbody></table></body><br/><br/>"
			if completeReport == "True" {

				for TaskName, value := range oldBuildDataTask[ResourceName] {
					//fmt.Printf("TaskName is:%s\n", TaskName)

					//ks := reflect.ValueOf(value).MapKeys()
					//fmt.Println(oldBuildDataTask[ResourceName][TaskName])
					val, _ := value.(map[string]interface{})

					if _, ok := newBuildDataTask[ResourceName][TaskName]; ok {
						p = p + fmt.Sprintf("<table style='backgound:#fff;border-collapse: collapse;' border = '1' cellpadding = '6'><tbody><tr><td colspan=5 style='text-align:center;background-color:#7388f9;color:white;'><b>Task Name : %s </b></td></tr><tr><th>Sub Task</th><th>Build# %s </th><th>Build# %s</th><th>Time Difference</th><th> %% Time Difference</th></tr> ", TaskName, oldBuildNum, newBuildNum)

						for key := range val {
							//fmt.Println(oldBuildDataTask[TaskName][key])

							svOldTask := oldBuildDataTask[ResourceName][TaskName].(map[string]interface{})[key].(string)
							svNewTask := newBuildDataTask[ResourceName][TaskName].(map[string]interface{})[key].(string)
							timeOldTask, _ := time.Parse(config.TimeFormat, svOldTask)
							timeNewTask, _ := time.Parse(config.TimeFormat, svNewTask)
							diffTask := timeNewTask.Sub(timeOldTask)

							if ((timeOldTask.Second() + (timeOldTask.Minute() * 60) + (timeOldTask.Hour() * 3600)) > 30) && ((timeNewTask.Second() + (timeNewTask.Minute() * 60) + (timeNewTask.Hour() * 3600)) > 30) {
								if diffTask <= 0 {
									percDiffTask := utils.CalcPerc(float64(diffTask.Seconds()), timeOldTask)

									p = p + "<tr style='background:#a7f392'><td>" + key + "</td><td>" + svOldTask + "</td><td>" + svNewTask + "</td><td>" + diffTask.String() + " </td><td>" + strconv.FormatFloat(percDiffTask, 'f', 2, 64) + " %</td></tr>"

								} else {

									percDiffTask := utils.CalcPerc(float64(diffTask.Seconds()), timeOldTask)
									p = p + "<tr style='background:#f99'><td>" + key + "</td><td>" + svOldTask + "</td><td>" + svNewTask + "</td><td>" + diffTask.String() + " </td><td>" + strconv.FormatFloat(percDiffTask, 'f', 2, 64) + " %</td></tr>"
								}
							}
						}
						p = p + "</tbody></table></body><br/><br/>"

					}

				}

			}

		}
		//fmt.Println(p)
		conf := utils.ReadConfig()
		p = p + "<b>Dashboard URL : </b>" + conf.DashboardURL
		utils.SendMail(p, subject, cc)
		//fmt.Println(p)
		utils.RespondWithJSON("Email Sent Successfully", w, r)
	}

}

func createJson(w http.ResponseWriter, r *http.Request) {

	Hostname := r.URL.Query().Get("hostname")
	Port := r.URL.Query().Get("port")
	ResourceName := r.URL.Query().Get("resourcename")
	Build := r.URL.Query().Get("build")
	Release := r.URL.Query().Get("release")
	User := r.URL.Query().Get("user")
	Pass := r.URL.Query().Get("pass")
	//regp := regexp.MustCompile("\\/\\/(.+)\\.informatica")
	//Host := regp.FindStringSubmatch(Hostname)
	flag := true
	URL := "http://" + Hostname + ".informatica.com:" + Port + "/ldmadmin/web.isp/login"
	JobURL := "http://" + Hostname + ".informatica.com:" + Port + "/ldmadmin/ldm.resources.core/execute"
	TaskURL := "http://" + Hostname + ".informatica.com:" + Port + "/ldmadmin/ldm.monitoring/jobprogress"
	var ldmHeader config.LDMHeader
	elasticJson := ""

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}
	var jobData config.JobResponse
	var resp *http.Response

	for len(jobData) < 1 {
		resp, client = utils.LDMLogin(URL, client, User, Pass)

		if resp == nil {
			flag = false
			utils.RespondWithText("Please check hostname or port ", w, r)
			break

		}

		// split cookie for next request
		var jsonData config.LDMResponse
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal([]byte(body), &jsonData)

		//check if login was not successful
		if jsonData.Error != "" {

			flag = false
			utils.RespondWithText("Incorrect username/pass ", w, r)
			break
		}

		if err != nil {
			fmt.Println("in err")
			fmt.Println(err)
		}
		fmt.Println(string(body))
		cookie := resp.Header.Get("Set-Cookie")
		newcookie := strings.Split(cookie, ";")
		jsession := strings.Split(newcookie[0], "=")
		resp.Header.Set("Set-Cookie", newcookie[0])

		//values = map[string]string{"name": "Profiling_5tables", "status": "true", "infaToken": jsonData.InfaToken}
		ldmHeader.InfaToken = jsonData.InfaToken
		ldmHeader.JsessionID = jsession[1]

		//failure := true
		rf := true
		jobData, client, rf = utils.GetJobData(resp, client, JobURL, ResourceName, ldmHeader)
		fmt.Println(len(jobData))
		if !rf {
			flag = false
			utils.RespondWithText("Please check Resource Name ", w, r)
			break
		}

	}

	//fmt.Println(jobData)
	if flag {

		for i := 0; i < len(jobData); i++ {
			//fmt.Println(jobData[i].Type)
			//sresp, _ := time.Parse(config.TimeFormat, resp).String()

			if jobData[i].Type != "Purge" {
				if jobData[i].Status != "SUCCESS" {

					utils.RespondWithText("Scanner Execution not completed or Failed! Please Check.", w, r)
					flag = false
					break
				}
			}
			taskResponseData := utils.GetTaskData(ldmHeader.InfaToken, jobData[i].ID, jobData[i].Taskid, TaskURL, resp, client)
			//fmt.Println(len(taskResponseData[0].Progress))

			if len(taskResponseData) > 0 {
				for j := 0; j < len(taskResponseData[0].Progress); j++ {

					//fmt.Println(taskResponseData[0].Progress[j].StartTime)
					//fmt.Println(taskResponseData[0].Progress[j].EndTime)
					startTime := time.Unix(0, taskResponseData[0].Progress[j].StartTime*int64(time.Millisecond))
					endTime := time.Unix(0, taskResponseData[0].Progress[j].EndTime*int64(time.Millisecond))
					diff := endTime.Sub(startTime)
					//fmt.Printf("%f", diff.Seconds()/
					//p := fmt.Sprintf("%02d:%02d:%02d", int64(diff.Hours()), int64(diff.Minutes()), int64(diff.Seconds()))
					sdiff := utils.SecToTime(int64(diff.Seconds()))
					//fmt.Println(sdiff)

					elasticJson, _ = sjson.Set(elasticJson, "TaskTimes."+jobData[i].Type+"."+taskResponseData[0].Progress[j].OperationName, sdiff)

				}
			}
			//fmt.Println(taskResponseData)

			elasticJson, _ = sjson.Set(elasticJson, "Times."+jobData[i].Type, jobData[i].ElapsedTime)
		}

		if flag {

			if len(jobData) > 1 && jobData[1].Type != "Purge" {
				endToEndTime := utils.EndToEndTime(jobData)
				elasticJson, _ = sjson.Set(elasticJson, "Times.End to End Execution Time", endToEndTime)
			}

			elasticJson, _ := sjson.Set(elasticJson, "ResourceName", ResourceName)
			elasticJson, _ = sjson.Set(elasticJson, "Hostname", Hostname)
			elasticJson, _ = sjson.Set(elasticJson, "Build", Build)
			elasticJson, _ = sjson.Set(elasticJson, "Release", Release)
			//fmt.Println(newcookie[0])
			//fmt.Println(elasticJson)

			//var t config.TimesResponse
			//json.Unmarshal([]byte(elasticJson), &t)
			//fmt.Println(t)
			rawIn := json.RawMessage(elasticJson)
			jsonBody, err := rawIn.MarshalJSON()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//fmt.Println(jsonBody)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonBody)
		}
	}
}

func htmlReport(w http.ResponseWriter, r *http.Request) {

	p := "../" + r.URL.Path
	//fmt.Println(r.URL.Path)

	http.ServeFile(w, r, p)

}
