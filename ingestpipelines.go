package main

import (
	"io"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

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

func setupIngestPipelines(client *elasticsearch.Client) (string, int, error) {
	return makeEsRequest(
		client,
		func(br io.Reader) (*esapi.Response, error) {
			return client.Ingest.PutPipeline("geoip", br)
		},
		IngestPipelineSetupReq{
			Description: "GeoIP ingest",
			Processors: []IngestPipelineProcessor{
				{
					GeoIP: GeoIPProcessor{
						Field: "sourceIP",
					},
				},
			},
		},
	)
}
