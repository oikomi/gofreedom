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

package utils

import (
	"fmt"
	"net"
	"net/http"
	"io"
)

func GetHostIP(host string) (ip string, err error) {
	ipaddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return ipaddr.String(), err
}

func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func CoreCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	var buffer [8192]byte
	buf := buffer[:]

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func CopyLink(dst, src io.ReadWriteCloser) {
	go func() {
		defer src.Close()
		CoreCopy(src, dst)
	}()
	defer dst.Close()
	CoreCopy(dst, src)
}

