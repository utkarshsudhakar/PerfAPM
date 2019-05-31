package config

type JobResponse []struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Taskid      string `json:"taskid"`
	Status      string `json:"status"`
	ElapsedTime string `json:"elapsedTime"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

type TaskResponse []struct {
	JobID    string `json:"jobId"`
	Progress []struct {
		//OperationNameLabel string `json:"operationNameLabel"`
		OperationName string `json:"operationName"`
		StartTime     int64  `json:"startTime"`
		EndTime       int64  `json:"endTime"`
	} `json:"progress"`
	TaskID string `json:"taskId"`
}

type ElasticsearchData struct {
	Hostname     string `json:"Hostname"`
	Release      string `json:"Release"`
	ResourceName string `json:"ResourceName"`
	Times        interface{}
	TaskTime     interface{}
}
