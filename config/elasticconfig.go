package config

const ElasticURL = "http://irl75cmgperf07:9200"

//const ElasticURL = "http://34.196.162.239:9200"

const ElasticSearchReportIndex = "regressioncomparison_14may"

const TimeFormat = "15:04:05"

type TimesResponse struct {
	Hostname     string      `json:"Hostname"`
	Release      string      `json:"Release"`
	Build        string      `json:"Build"`
	ResourceName string      `json:"ResourceName"`
	Times        interface{} `json:"Times" bson:"Times"`
	TaskTimes    interface{} `json:"TaskTimes" bson:"TaskTimes"`
}