package main

import (
	"net/http"
)

type DataViewInfo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type DataViewReq struct {
	DataView DataViewInfo `json:"data_view"`
}

func dataViewExists(client *http.Client, view DataViewInfo) bool {
	_, code, _ := makeKibanaRequest(
		client,
		"GET",
		"/data_views/data_view/"+view.Id,
		"",
	)
	if code == 200 {
		return true
	}
	return false
}

func setupDataView(client *http.Client, dataview DataViewInfo) (string, int, error) {
	if dataViewExists(client, dataview) {
		return "Data view already exists.", 200, nil
	}
	res, code, err := makeKibanaRequest(
		client,
		"POST",
		"/data_views/data_view",
		DataViewReq{
			DataView: dataview,
		},
	)
	return string(res), code, err
}
