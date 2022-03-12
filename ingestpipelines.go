package main

import elasticsearch "github.com/elastic/go-elasticsearch/v8"

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
