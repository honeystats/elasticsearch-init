package main

import (
	"io"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func indexExists(client *elasticsearch.Client, index string) bool {
	_, code, _ := makeEsRequest(
		client,
		func(br io.Reader) (*esapi.Response, error) {
			return client.Indices.Get([]string{index})
		},
		locationMappingReq,
	)
	return code == 200
}

func createIndex(client *elasticsearch.Client, index string) (string, int, error) {
	if indexExists(client, index) {
		return "Index already exists.", 200, nil
	}
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
	return res, 200, nil
}
