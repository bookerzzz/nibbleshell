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
	"os"
	"text/template"
)

// Nibbleshell is the primary struct of the program. It holds onto the
// configuration, the HTTP server, and all the routes.
type Nibbleshell struct {
	Pid    int
	Config *Config
	Routes []*Route
	Server *Server
}

// NewWithConfig creates a new instance from an instance of Config.
func NewWithConfig(config *Config) (*Nibbleshell, error) {
	routes := make([]*Route, 0, len(config.RouteConfigs))
	for _, routeConfig := range config.RouteConfigs {
		route, err := NewRouteWithConfig(routeConfig, config.StatterConfig)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	return &Nibbleshell{
		Pid:    os.Getpid(),
		Config: config,
		Routes: routes,
		Server: NewServerWithConfigAndRoutes(config.ServerConfig, routes),
	}, nil
}

// Run starts the HTTP server of the service.
func (h *Nibbleshell) Run() error {
	tmpl, err := template.New("start").Parse(StartupTemplateString)
	if err != nil {
		return err
	}
	err = tmpl.Execute(os.Stdout, h)
	if err != nil {
		return err
	}

	return h.Server.ListenAndServe()
}
