package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type DataViewInfo struct {
	Title string `json:"title"`
}

type DataViewReq struct {
	DataView DataViewInfo `json:"data_view"`
}

type DataViewErr struct {
	Message string `json:"message"`
}

func setupDataView(client *http.Client, view string) (string, int, error) {
	res, code, err := makeKibanaRequest(
		client,
		"POST",
		"/data_views/data_view",
		DataViewReq{
			DataView: DataViewInfo{
				Title: view,
			},
		},
	)
	if code == 400 {
		var dve DataViewErr
		json.Unmarshal([]byte(res), &dve)
		if strings.Contains(dve.Message, "Duplicate index pattern:") {
			return "Data view already exists.", 200, nil
		}
	}
	return res, code, err
}
