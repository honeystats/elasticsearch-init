package main

import (
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

var ELASTICSEARCH_URL string

func requireEnv(key string) string {
	val, isSet := os.LookupEnv(key)
	if !isSet {
		logrus.Fatalf("Missing environment variable: [$%s]\n", key)
	}
	return val
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})
	ELASTICSEARCH_URL = requireEnv("ELASTICSEARCH_URL")
}

func main() {
	logrus.Info("Starting ES setup...")
	cfg := elasticsearch.Config{
		Addresses: []string{
			ELASTICSEARCH_URL,
		},
	}
	_, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up ES client")
	}
	logrus.Info("Elasticsearch client initialized successfully.")
}
