// Copyright (c) 2014 Oyster
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

// Halfshell is the primary struct of the program. It holds onto the
// configuration, the HTTP server, and all the routes.
type Halfshell struct {
	Pid    int
	Config *Config
	Routes []*Route
	Server *Server
	Logger *Logger
}

// NewWithConfig creates a new Halfshell instance from an instance of Config.
func NewWithConfig(config *Config) *Halfshell {
	routes := make([]*Route, 0, len(config.RouteConfigs))
	for _, routeConfig := range config.RouteConfigs {
		routes = append(routes, NewRouteWithConfig(routeConfig, config.StatterConfig))
	}

	return &Halfshell{
		Pid:    os.Getpid(),
		Config: config,
		Routes: routes,
		Server: NewServerWithConfigAndRoutes(config.ServerConfig, routes),
		Logger: NewLogger("main"),
	}
}

// Run starts the Halfshell program. Performs global (de)initialization, and
// starts the HTTP server.
func (h *Halfshell) Run() {
	var tmpl, _ = template.New("start").Parse(StartupTemplateString)
	_ = tmpl.Execute(os.Stdout, h)

	imagick.Initialize()
	defer imagick.Terminate()

	h.Server.ListenAndServe()
}
