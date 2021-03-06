/*
Copyright 2020 The Netease Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type StringFlagSetterFunc func(string) (string, error)

// StringFlagPutHandler wraps an http Handler to set string type flag.
func StringFlagPutHandler(setter StringFlagSetterFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch {
		case req.Method == "PUT":
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				writePlainText(http.StatusBadRequest, "error reading request body: "+err.Error(), w)
				return
			}
			defer req.Body.Close()
			response, err := setter(string(body))
			if err != nil {
				writePlainText(http.StatusBadRequest, err.Error(), w)
				return
			}
			writePlainText(http.StatusOK, response, w)
			return
		default:
			writePlainText(http.StatusNotAcceptable, "unsupported http method", w)
			return
		}
	})
}

// writePlainText renders a simple string response.
func writePlainText(statusCode int, text string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, text)
}
