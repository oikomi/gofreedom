//
// Copyright 2014 Hong Miao. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpproxy
 
import (
	"net/http"
	"fmt"
	"io"
	"time"
	"../config"
	//"../httplib"
)
type Handler int8 
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (h *Handler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func HTTPProxyServer(cfg config.Config) {
	port := cfg.Listen
	s := &http.Server {
		Addr:    port,
		Handler: new(Handler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}