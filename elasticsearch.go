package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func parseEsRes(resObj *esapi.Response, resErr error) ([]byte, int, error) {
	if resErr != nil {
		return []byte{}, -1, resErr
	}
	defer resObj.Body.Close()

	out, err := ioutil.ReadAll(resObj.Body)
	if err != nil {
		return nil, -1, errors.New(fmt.Sprintf("Error reading ES response: %v", err))
	}

	if resObj.IsError() {
		return nil, resObj.StatusCode, errors.New(fmt.Sprintf("Status text: [%s], Response text: [%s]", resObj.Status(), string(out)))
	}

	return out, resObj.StatusCode, nil
}

func setupEsReq(body interface{}) (io.Reader, error) {
	byteArr, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error marshalling ES request: %v", err))
	}
	return bytes.NewReader(byteArr), nil
}

func makeEsRequest(client *elasticsearch.Client, req func(br io.Reader) (*esapi.Response, error), body interface{}) (string, int, error) {
	bodyReader, err := setupEsReq(body)
	if err != nil {
		return "", -1, err
	}
	txt, status, err := parseEsRes(req(bodyReader))
	return string(txt), status, err
}
