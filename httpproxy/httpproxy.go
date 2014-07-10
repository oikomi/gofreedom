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
	"log"
	//"time"
	"net"
	"os"
	"net/http/httputil"
	"strings"
	"../config"
	"../utils"
	"../httplib"
)

const BUF_SIZE = 65535

type HTTPProxy struct {
	//transport http.Transport
	//mux       *http.ServeMux
	logger *log.Logger
}

func NewProxy(logger *log.Logger) (p *HTTPProxy) {
	p = &HTTPProxy {
		logger : logger,
	}
	return p
}

func upstream(tcpaddr string, req_header []byte) []byte {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpaddr)
	if err != nil {
		fmt.Println(err.Error())
		return []byte("")
	}
	utils.CheckError(err, "ResolveTCPAddr")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	utils.CheckError(err, "DialTCP")

	_, err = conn.Write(req_header)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return []byte("")
	}

	buf := make([]byte, BUF_SIZE)
	for {
		lenght, err := conn.Read(buf)
		if utils.CheckError(err, "Connection") == false {
			conn.Close()
			fmt.Println("Server is dead ...ByeBye")
			os.Exit(0)
		}
		fmt.Println("recive upstream response:  ")
		//fmt.Println(string(buf[0:lenght]))
		return buf[0:lenght]
	}
	return []byte("")
}

func (p *HTTPProxy)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.RequestURI)
	p.logger.Printf("req from %s visit %s \n" , r.RemoteAddr, r.RequestURI)
	
	if r.Method == "CONNECT" {
		p.Connect(w, r)
		return
	}
	tcpaddr, err := utils.GetHostIP(r.Host)
	if tcpaddr == "" {
		return 
	}
	dump , err := httputil.DumpRequest(r, false)
	//fmt.Println(string(dump))
	if err != nil {
		fmt.Println(string(dump))
		return
	}

	//resp := upstream(tcpaddr + ":80", dump)
	
	//fmt.Println(string(resp))

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		w.WriteHeader(httplib.StatusNotFound)
		return
	}
	defer resp.Body.Close()
	utils.CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_ ,err = io.Copy(w, resp.Body)
	if err != nil {
		w.WriteHeader(httplib.StatusInternalServerError)
		return
	}
}


func (p *HTTPProxy) Connect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect")
	hij, ok := w.(http.Hijacker)
	if !ok {
		p.logger.Fatal("httpserver does not support hijacking")
		return
	}
	srcconn, _, err := hij.Hijack()
	if err != nil {
		p.logger.Fatal("Cannot hijack connection ", err)
		return
	}
	defer srcconn.Close()

	host := r.URL.Host
	if !strings.Contains(host, ":") {
		host += ":80"
	}
	dstconn, err := net.Dial("tcp", host)
	if err != nil {
		p.logger.Fatal("dial failed:", err)
		srcconn.Write([]byte("HTTP/1.0 502 OK\r\n\r\n"))
		return
	}
	//srcconn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	srcconn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	utils.CopyLink(srcconn, dstconn)
	return
}

func HTTPProxyServer(cfg *config.Config, logger *log.Logger) {
	err := http.ListenAndServe(cfg.Listen, NewProxy(logger))
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		logger.Fatal("ListenAndServe: ", err)
		return
	}
}