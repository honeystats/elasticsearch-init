package main

import (
	"io"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var HoneystatsIndices = []string{"ssh_data", "web_data"}

func createIndices(client *elasticsearch.Client) (string, int, error) {
	for _, index := range HoneystatsIndices {
		res, code, err := makeEsRequest(
			client,
			func(br io.Reader) (*esapi.Response, error) {
				return client.Indices.Create(index)
			},
			locationMappingReq,
		)
		if err != nil {
			return res, code, err
		}
	}
	return "Indices created.", 200, nil
}
