package rest

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

type HandleFunc func(r *http.Request) (*HTTPData, error)

type HTTPData struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func readPostData(r *http.Request, target interface{}) error {
	req, _ := httputil.DumpRequest(r, true)
	logrus.Infoln(string(req))
	if r.Body == nil{
		logrus.Warn("Body nil")
	}
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&target)
}