package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func parseKibanaRes(resObj *http.Response, resErr error) ([]byte, int, error) {
	if resErr != nil {
		return []byte{}, -1, resErr
	}
	defer resObj.Body.Close()

	out, err := ioutil.ReadAll(resObj.Body)
	if err != nil {
		return nil, -1, errors.New(fmt.Sprintf("Error reading Kibana response: %v", err))
	}

	if resObj.StatusCode >= 400 {
		return out, resObj.StatusCode, errors.New(fmt.Sprintf("Status text: [%s]", resObj.Status))
	}

	return out, resObj.StatusCode, nil
}

func setupKibanaReq(body interface{}) (io.Reader, error) {
	byteArr, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error marshalling Kibana request: %v", err))
	}
	return bytes.NewReader(byteArr), nil
}

func makeKibanaRequest(client *http.Client, method string, path string, body interface{}) (string, int, error) {
	bodyReader, err := setupKibanaReq(body)
	if err != nil {
		return "", -1, err
	}
	reqPath := KIBANA_API_URL + path
	req, err := http.NewRequest(method, reqPath, bodyReader)
	req.SetBasicAuth("elastic", ELASTIC_PASSWORD)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("kbn-xsrf", "true")
	if err != nil {
		return "", -1, err
	}
	txt, status, err := parseKibanaRes(client.Do(req))
	return string(txt), status, err
}
