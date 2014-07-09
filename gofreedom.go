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

package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"./glog"
	"./config"
	"./httpproxy"
)


/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

const char* build_time(void) {
	static const char* psz_build_time = "["__DATE__ " " __TIME__ "]";
	return psz_build_time;
}
*/
import "C"

var (
	buildTime = C.GoString(C.build_time())
)

func BuildTime() string {
	return buildTime
}

const VERSION string = "0.01"

const BUF_SIZE = 65535

//todo
func enc(raw string) {

}

func dec(raw string) {

}

func checkError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + "  " + err.Error())
		return false
	}
	return true
}

func upstream(tcpaddr string, req_header string) []byte {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpaddr)
	checkError(err, "ResolveTCPAddr")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "DialTCP")

	_, err = conn.Write([]byte(req_header))
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
	}

	buf := make([]byte, BUF_SIZE)
	for {
		lenght, err := conn.Read(buf)
		if checkError(err, "Connection") == false {
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

func Handler(conn net.Conn, messages chan string) {
	fmt.Println("connection is connected from ...", conn.RemoteAddr().String())
	buf := make([]byte, BUF_SIZE)
	for {
		lenght, err := conn.Read(buf)
		if checkError(err, "Connection") == false {
			conn.Close()
			break
		}
		if lenght > 0 {
			buf[lenght] = 0
		}

		reciveStr := string(buf[0:lenght])
		messages <- reciveStr
	}
}

func getHostIP(buf string) string {
	for _, s := range strings.Split(buf, "\n") {
		index := strings.Index(s, "Host:")
		if index != -1 {
			host := s[6 : len(s)-1]
			fmt.Println(host)
			fmt.Println(reflect.TypeOf(host))
			ipaddr, err := net.ResolveIPAddr("ip", string(host))
			if err != nil {
				fmt.Println(err)
			}

			return ipaddr.String()
		}
	}

	return ""
}

func forwardHandler(conns *map[string]net.Conn, messages chan string, conn net.Conn) {
	for {
		req_header := <-messages

		ipaddr := getHostIP(req_header)
		
		if ipaddr != "" {
			conn.Write(upstream(ipaddr+":80", req_header))
		}

		//for key, value := range *conns {
		//	fmt.Println("connection is connected from ...", key)
		//	_, err := value.Write([]byte(req_header))
		//	if err != nil {
		//		fmt.Println(err.Error())
		//		//delete(*conns, key)
		//	}
		//}
	}
}

func StartServer(port string) {
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")
	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")
	conns := make(map[string]net.Conn)
	messages := make(chan string, 100)

	for {
		fmt.Printf("Listening on port %s ...\n", port)
		conn, err := l.Accept()
		checkError(err, "Accept")
		fmt.Println("Accepting ...")
		conns[conn.RemoteAddr().String()] = conn
		go Handler(conn, messages)
		go forwardHandler(&conns, messages, conn)
	}
}
//////////////////////////

func StartHttpServer(cfg config.Config) {
	httpproxy.HTTPProxyServer(cfg)
}

func usage() {
	fmt.Printf("Usage : gofreedom port  \n")
}

func version() {
	fmt.Printf("gofreedom version %s Copyright (c) 2014 Harold Miao (miaohonghit@gmail.com)  \n", VERSION)
}

func main() {
	version()
	fmt.Printf("built on %s\n", BuildTime())
	if len(os.Args) != 2 {
		glog.Error("Usage : gofreedom config_file")
		os.Exit(0)
	}
	
	cfg, err := config.LoadConfig(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	config.DumpConfig(cfg)
	StartHttpServer(cfg)
	//StartServer(os.Args[1])
}
