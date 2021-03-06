// Copyright (c) 2014 Oyster
// Copyright (c) 2015 Hotel Booker B.V.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package nibbleshell

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type Statter interface {
	RegisterRequest(*ResponseWriter, *Request)
}

type statsdStatter struct {
	conn     *net.UDPConn
	addr     *net.UDPAddr
	Name     string
	Hostname string
	Enabled  bool
}

func NewStatterWithConfig(routeConfig *RouteConfig, statterConfig *StatterConfig) (Statter, error) {
	s := &statsdStatter{}

	if statterConfig.Enabled {
		var err error
		s.Hostname, err = os.Hostname()
		if err != nil {
			return s, err
		}

		s.addr, err = net.ResolveUDPAddr(
			"udp", fmt.Sprintf("%s:%d", statterConfig.Host, statterConfig.Port))
		if err != nil {
			return s, err
		}

		s.conn, err = net.DialUDP("udp", nil, s.addr)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

func (s *statsdStatter) RegisterRequest(w *ResponseWriter, r *Request) {
	if !s.Enabled {
		return
	}

	now := time.Now()

	status := "success"
	if w.Status != http.StatusOK {
		status = "failure"
	}

	s.count(fmt.Sprintf("http.status.%d", w.Status))
	s.count(fmt.Sprintf("image_resized.%s", status))
	s.count(fmt.Sprintf("image_resized_%s.%s", r.ProcessorOptions.String(), status))

	if status == "success" {
		durationInMs := (now.UnixNano() - r.Timestamp.UnixNano()) / 1000000
		s.time("image_resized", durationInMs)
		s.time(fmt.Sprintf("image_resized_%s", r.ProcessorOptions.String()), durationInMs)
	}
}

func (s *statsdStatter) count(stat string) {
	stat = fmt.Sprintf("%s.nibbleshell.%s.%s", s.Hostname, s.Name, stat)
	s.send(stat, "1|c")
}

func (s *statsdStatter) time(stat string, time int64) {
	stat = fmt.Sprintf("%s.nibbleshell.%s.%s", s.Hostname, s.Name, stat)
	s.send(stat, fmt.Sprintf("%d|ms", time))
}

func (s *statsdStatter) send(stat string, value string) error {
	data := fmt.Sprintf("%s:%s", stat, value)
	_, err := s.conn.Write([]byte(data))
	return err
}
