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
	"net/http/pprof"
	runtime_pprof "runtime/pprof"
	"strings"

	"k8s.io/klog/v2"
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

// EnableProfiling enables golang profiling.
func EnableProfiling() {
	for _, profile := range runtime_pprof.Profiles() {
		name := profile.Name()
		handler := pprof.Handler(name)
		addPath(name, handler)
	}

	// static profiles as listed in net/http/pprof/pprof.go:init()
	addPath("cmdline", http.HandlerFunc(pprof.Cmdline))
	addPath("profile", http.HandlerFunc(pprof.Profile))
	addPath("symbol", http.HandlerFunc(pprof.Symbol))
	addPath("trace", http.HandlerFunc(pprof.Trace))
}

func addPath(name string, handler http.Handler) {
	http.Handle(name, handler)
	klog.V(4).Infof("DEBUG: registered profiling handler on /debug/pprof/%s", name)
}

func HttpGet(reqURL string, queryMap map[string]string) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, nil, err
	}

	// set http header
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Accept", "application/json")
	req.Header = header

	// TODO: now the snapshot server can not recognize URL encode
	req.URL.RawQuery = buildRawQuery(queryMap)

	klog.V(6).Infof("Request: %+v", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, nil, err
	}
	klog.V(6).Infof("Response: %+v", *resp)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	klog.V(7).Infof("Response data: %+v", string(data))
	return resp.StatusCode, data, nil
}

func buildRawQuery(queryMap map[string]string) string {
	var buf strings.Builder
	for k, v := range queryMap {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}
	return buf.String()
}
