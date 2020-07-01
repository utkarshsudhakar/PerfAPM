package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic"
)

//GetBuildData ...
func CheckDuplicate(buildNum string, release string, Hostname string, ResourceName string) bool {

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

	filterByBuildQuery := elastic.NewTermQuery("Build", buildNum)
	filterByReleaseQuery := elastic.NewTermQuery("Release", release)
	//searchQuery := elastic.NewTermQuery("Hostname", Hostname).NewRegexpQuery()
	searchQuery := elastic.NewRegexpQuery("Hostname", Hostname)
	resourceNameQuery := elastic.NewTermQuery("ResourceName.keyword", ResourceName)
	filterQuery := elastic.NewBoolQuery().Must(filterByReleaseQuery).Must(filterByBuildQuery).Must(searchQuery).Must(resourceNameQuery)

	//for filter based on last build num use aggregation max with release

	SearchResult, err := client.Search().
		Index(conf.ElasticSearchReportIndex). // search in index "testutkarsh"
		Query(filterQuery).
		From(0).Size(1000).
		Pretty(true).
		Do(context.Background())

	if err != nil {
		panic(err)
	}
	if SearchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d hits\n", SearchResult.Hits.TotalHits)

		return true

	}

	// No hits
	msg := fmt.Sprintf("Found no hits for resourceName: %s", ResourceName)

	fmt.Println(msg)

	return false

}
