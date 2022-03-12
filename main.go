package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up ES client")
	}
	logrus.Info("Elasticsearch client initialized successfully.")

	res, err := setupIngestPipelines(client)
	if err != nil {
		logrus.WithError(err).Fatalln("Error setting up ingest pipeline.")
	}
	logrus.WithField("result", res).Infoln("Succesfully set up ingest pipeline.")
}

func parseEsRes(resObj *esapi.Response, resErr error) ([]byte, error) {
	defer resObj.Body.Close()

	var out []byte
	_, err := resObj.Body.Read(out)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading ES response: %v", err))
	}

	if resObj.IsError() {
		return nil, errors.New(fmt.Sprintf("Status text: [%s], Response text: [%s]", resObj.Status(), string(out)))
	}

	return out, nil
}

func setupEsReq(body interface{}) (io.Reader, error) {
	byteArr, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error marshalling ES request: %v", err))
	}
	return bytes.NewReader(byteArr), nil
}

type IngestPipelineSetupReq struct {
	Description string                    `json:"description"`
	Processors  []IngestPipelineProcessor `json:"processors"`
}

type IngestPipelineProcessor struct {
	GeoIP GeoIPProcessor `json:"geoip"`
}
type GeoIPProcessor struct {
	Field string `json:"field"`
}

func setupIngestPipelines(client *elasticsearch.Client) (string, error) {
	body := IngestPipelineSetupReq{
		Description: "GeoIP ingest",
		Processors: []IngestPipelineProcessor{
			{
				GeoIP: GeoIPProcessor{
					Field: "sourceIP",
				},
			},
		},
	}
	bodyReader, err := setupEsReq(body)
	if err != nil {
		return "", err
	}

	txt, err := parseEsRes(client.Ingest.PutPipeline("geoip", bodyReader))
	if err != nil {
		return "", err
	}

	return string(txt), nil
}
