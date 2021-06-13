package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eucatur/go-toolbox/log"
	"github.com/jmoiron/sqlx/types"
	"github.com/parnurzeal/gorequest"
	"time"
)

func logFile(response gorequest.Response, tempStruct interface{}, body []byte, errs []error) {
	var sent, received interface{}

	if response.Request.Body != nil {
		func() {
			sentBody, err := response.Request.GetBody()
			if err != nil {
				fmt.Println(err)
				return
			}

			buf := new(bytes.Buffer)

			_, err = buf.ReadFrom(sentBody)
			if err != nil {
				fmt.Println(err)
				return
			}

			sentBytes := buf.Bytes()
			if json.Valid(sentBytes) {
				sent = types.JSONText(sentBytes)
			} else {
				sent = string(sentBytes)
			}
		}()
	}

	if json.Valid(body) {
		received = types.JSONText(body)
	} else {
		received = string(body)
	}

	logBody := struct {
		Method   string      `json:"method"`
		Scheme   string      `json:"scheme"`
		Host     string      `json:"host"`
		URL      string      `json:"url"`
		Sent     interface{} `json:"sent"`
		Status   int         `json:"status"`
		Received interface{} `json:"received"`
	}{
		Method:   response.Request.Method,
		Scheme:   response.Request.URL.Scheme,
		Host:     response.Request.URL.Hostname(),
		URL:      response.Request.URL.RequestURI(),
		Sent:     sent,
		Status:   response.StatusCode,
		Received: received,
	}

	logBytes, err := json.Marshal(logBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = log.File(time.Now().Format("clients/api-srvp/2006/01/02/15h.log"), string(logBytes))
	if err != nil {
		fmt.Println(err)
	}

	return
}
