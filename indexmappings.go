package main

import (
	"io"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type IndexMappingBaseField struct {
	Type string `json:"type"`
}

type IndexMappingSubProperty map[string]IndexMappingBaseField

type IndexMappingProperty struct {
	Properties IndexMappingSubProperty `json:"properties"`
}

type IndexMappingProperties map[string]IndexMappingProperty

type IndicesPutMappingReq struct {
	Properties IndexMappingProperties `json:"properties"`
}

var locationMappingReq = IndicesPutMappingReq{
	Properties: map[string]IndexMappingProperty{
		"location": {
			Properties: map[string]IndexMappingBaseField{
				"lat": {
					Type: "float",
				},
				"long": {
					Type: "float",
				},
			},
		},
	},
}

func setupLocationMapping(client *elasticsearch.Client) (string, int, error) {
	return makeEsRequest(
		client,
		func(br io.Reader) (*esapi.Response, error) {
			return client.Indices.PutMapping(HoneystatsIndices, br)
		},
		locationMappingReq,
	)
}
