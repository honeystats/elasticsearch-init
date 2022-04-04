package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

func setupDashboard(client *http.Client, filename string) (string, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "Error opening file", -1, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	stat, err := file.Stat()
	if err != nil {
		return "Error stat-ing file", -1, err
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(filename)))
	h.Set("Content-Type", "application/octet-stream")
	h.Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	part, err := writer.CreatePart(h)
	if err != nil {
		return "Error creating form file", -1, err
	}

	io.Copy(part, file)
	writer.Close()

	res, code, err := makeKibanaRequest(
		client,
		"POST",
		"/saved_objects/_import?overwrite=true",
		ExtraHeaders{
			"Content-Type": writer.FormDataContentType(),
		},
		body,
	)
	return string(res), code, err
}

func importDashboardsFromDir(client *http.Client) error {
	filenames, err := filepath.Glob("./dashboards/*.ndjson")
	if err != nil {
		return err
	}
	for _, filename := range filenames {
		res, code, err := setupDashboard(client, filename)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"res":      res,
				"code":     code,
				"filename": filename,
			}).Fatalln("Error setting up dashboard.")
		}
		logrus.WithFields(logrus.Fields{
			"res":      res,
			"code":     code,
			"filename": filename,
		}).Infoln("Successfully set up dashboard.")
	}
	return nil
}
