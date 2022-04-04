package main

import (
	"os"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

var ELASTIC_USERNAME string
var ELASTIC_PASSWORD string
var KIBANA_SYSTEM_PASSWORD string
var ELASTIC_URL string

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
	ELASTIC_USERNAME = requireEnv("ELASTIC_USERNAME")
	ELASTIC_PASSWORD = requireEnv("ELASTIC_PASSWORD")
	KIBANA_SYSTEM_PASSWORD = requireEnv("KIBANA_SYSTEM_PASSWORD")
	ELASTIC_URL = requireEnv("ELASTICSEARCH_URL")
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
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up ES client")
	}
	logrus.Info("Elasticsearch client initialized successfully.")

	res, status, err := setupIngestPipelines(client)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up ingest pipeline.")
	}
	logrus.WithFields(logrus.Fields{
		"result":      res,
		"status_code": status,
	}).Infoln("Succesfully set up ingest pipeline.")

	for _, index := range HoneystatsIndices {
		res, status, err = createIndex(client, index)
		if err != nil {
			logrus.WithError(err).Fatalln("Error creating indices.")
		}
		logrus.WithFields(logrus.Fields{
			"index":       index,
			"result":      res,
			"status_code": status,
		}).Infoln("Succesfully set up index.")
	}

	res, status, err = setupLocationMapping(client)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up location mapping.")
	}
	logrus.WithFields(logrus.Fields{
		"result":      res,
		"status_code": status,
	}).Infoln("Succesfully set up location mapping.")
}
