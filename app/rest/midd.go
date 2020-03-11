package rest

import (
	"net/http"
	"reflect"

	"encoding/json"
	"github.com/sirupsen/logrus"
)

var (
	codes = map[string]int{
		"InternalServerError": http.StatusInternalServerError,
		"BadRequest":          http.StatusBadRequest,
	}
)

func (ra *restApplication) handleUseCase(fn HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		httpData, err := fn(r)
		write(w, httpData, err)
	}
}

func write(w http.ResponseWriter, httpData *HTTPData, err error) {
	var httpCode int
	if err != nil {
		httpCode = codes[err.Error()]
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}
		w.WriteHeader(httpCode)
	}
	if httpCode < http.StatusInternalServerError {
		writeData(httpData, w)
	}
}

func writeData(httpData *HTTPData, writer http.ResponseWriter) {
	if httpData != nil && (!isIFaceNil(httpData.Data) || !isIFaceNil(httpData.Message)) {
		if isIFaceNil(httpData.Data) {
			httpData.Data = nil
		}
		bytes, _ := json.Marshal(httpData)
		if _, err := writer.Write(bytes); err != nil {
			logrus.Fatalln("err write response", err.Error())
			return
		}
	}
}

func isIFaceNil(c interface{}) bool {
	return c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil())
}
