package main

import (
	"net/http"
	"os"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

var ELASTIC_USERNAME string
var ELASTIC_PASSWORD string
var ELASTIC_URL string
var KIBANA_API_URL string

func requireEnv(key string) string {
	val, isSet := os.LookupEnv(key)
	if !isSet {
		logrus.Fatalf("Missing environment variable: $%s", key)
	}
	return val
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	ELASTIC_USERNAME = requireEnv("ELASTIC_USERNAME")
	ELASTIC_PASSWORD = requireEnv("ELASTIC_PASSWORD")
	ELASTIC_URL = requireEnv("ELASTICSEARCH_URL")
	KIBANA_API_URL = requireEnv("KIBANA_API_URL")
}

var HoneystatsIndices = []string{"honeystats_ssh_data", "honeystats_web_data"}

func main() {
	logrus.Info("Starting ES setup...")
	cfg := elasticsearch.Config{
		Addresses: []string{
			ELASTIC_URL,
		},
		Username:   ELASTIC_USERNAME,
		Password:   ELASTIC_PASSWORD,
		MaxRetries: 1000,
		RetryBackoff: func(attempt int) time.Duration {
			return time.Second * 10
		},
	}
	elasticClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up ES client")
	}
	logrus.Info("Elasticsearch client initialized successfully.")

	kibanaClient := new(http.Client)

	res, status, err := setupIngestPipelines(elasticClient)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up ingest pipeline.")
	}
	logrus.WithFields(logrus.Fields{
		"result":      res,
		"status_code": status,
	}).Infoln("Successfully set up ingest pipeline.")

	for _, index := range HoneystatsIndices {
		res, status, err = createIndex(elasticClient, index)
		if err != nil {
			logrus.WithError(err).Fatalln("Error creating indices.")
		}
		logrus.WithFields(logrus.Fields{
			"index":       index,
			"result":      res,
			"status_code": status,
		}).Infoln("Successfully set up index.")
	}

	res, status, err = setupLocationMapping(elasticClient)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up location mapping.")
	}
	logrus.WithFields(logrus.Fields{
		"result":      res,
		"status_code": status,
	}).Infoln("Successfully set up location mapping.")

	var dataviews []DataViewInfo = []DataViewInfo{
		{
			Id:    "honeystats-all",
			Title: "honeystats_*",
		},
	}
	for _, index := range HoneystatsIndices {
		dataviews = append(dataviews, DataViewInfo{
			Id:    index,
			Title: index,
		})
	}

	for _, dv := range dataviews {
		res, status, err := setupDataView(kibanaClient, dv)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"id":     dv.Id,
				"title":  dv.Title,
				"status": status,
				"res":    res,
			}).Fatalln("Error setting up data view.")
		}
		logrus.WithFields(logrus.Fields{
			"id":     dv.Id,
			"title":  dv.Title,
			"status": status,
			"res":    res,
		}).Infoln("Successfully set up data view.")
	}
}
