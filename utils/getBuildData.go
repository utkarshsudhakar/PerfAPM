package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/utkarshsudhakar/PerfAPM/config"
)

//GetBuildData ...
func GetBuildData(buildNum string) map[string]map[string]interface{} {

	//Hostname := "irl62dqd07"
	conf := ReadConfig()

	client, err := elastic.NewClient(
		elastic.SetURL(conf.ElasticURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {

	}

	//ping to check connectivity

	info, code, err := client.Ping(conf.ElasticURL).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Get doc for the specific buildnumber
	filterByBuildQuery := elastic.NewTermQuery("Build", buildNum)
	//searchQuery := elastic.NewRegexpQuery("Hostname", "sql.*")
	//labelQuery := elastic.NewFilterAggregation
	//dataQuery := elastic.NewBoolQuery().Must(labelQuery).Must(filterByBuildQuery).
	//newquery := elastic.NewBoolQuery().Must(searchQuery).Must(filterByBuildQuery)

	//for filter based on last build num use aggregation max with release

	SearchResult, err := client.Search().
		Index(conf.ElasticSearchReportIndex). // search in index "testutkarsh"
		Query(filterByBuildQuery).
		Pretty(true).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	if SearchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d hits\n", SearchResult.Hits.TotalHits)

		var t config.TimesResponse
		var myMap map[string]interface{}
		newMap := make(map[string]map[string]interface{})

		for _, hit := range SearchResult.Hits.Hits {

			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
				fmt.Printf("%s", "Error")
			}

			if t.Times != nil {
				key := strings.Split(string(t.ResourceName), "_")
				myMap = t.Times.(map[string]interface{})

				newMap[key[0]] = myMap
			}

		}

		return newMap

	}

	// No hits
	msg := fmt.Sprintf("Found no hits for buildNumber: %s", buildNum)

	fmt.Println(msg)

	return nil

}
