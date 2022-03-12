package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

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
