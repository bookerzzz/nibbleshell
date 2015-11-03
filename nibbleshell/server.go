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
	"time"

	log "github.com/Sirupsen/logrus"
)

type Server struct {
	*http.Server
	Routes []*Route
}

func NewServerWithConfigAndRoutes(config *ServerConfig, routes []*Route) *Server {
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Port),
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server := &Server{httpServer, routes}
	httpServer.Handler = server
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hw := s.NewResponseWriter(w)
	hr, err := s.NewRequest(r)
	if err != nil {
		log.WithError(err).Warning("Error processing request")
		hw.WriteError("Bad Request", http.StatusBadRequest)
		return
	}
	defer s.LogRequest(hw, hr)
	switch {
	case "/healthcheck" == hr.URL.Path || "/health" == hr.URL.Path:
		hw.Write([]byte("OK"))
	default:
		s.ImageRequestHandler(hw, hr)
	}
}

func (s *Server) ImageRequestHandler(w *ResponseWriter, r *Request) {
	log := log.WithField("handler", "image")
	if r.Route == nil {
		w.WriteError(fmt.Sprintf("No route available to handle request: %v",
			r.URL.Path), http.StatusNotFound)
		return
	}

	defer func() { go r.Route.Statter.RegisterRequest(w, r) }()

	log.Info(fmt.Sprintf("Handling request for image %s with processor options %q",
		r.SourceOptions.Path, r.ProcessorOptions))

	image, err := r.Route.Source.GetImage(r.SourceOptions)
	if err != nil {
		w.WriteError("Not Found", http.StatusNotFound)
		return
	}

	// here image gets overriden with the processed version
	image, err = r.ProcessorOptions.ProcessImage(image)
	if err != nil {
		log.WithError(err).Warning("Error processing image data")
		w.WriteError("Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.SetHeader("Cache-Control", r.Route.CacheControl)
	w.WriteImage(image)
}

func (s *Server) LogRequest(w *ResponseWriter, r *Request) {
	logFormat := "%s - - [%s] \"%s %s %s\" %d %d\n"
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	fmt.Printf(logFormat, host, r.Timestamp.Format("02/Jan/2006:15:04:05 -0700"),
		r.Method, r.URL.RequestURI(), r.Proto, w.Status, w.Size)
}

type Request struct {
	*http.Request
	Timestamp        time.Time
	Route            *Route
	SourceOptions    *ImageSourceOptions
	ProcessorOptions *ImageProcessorOptions
}

func (s *Server) NewRequest(r *http.Request) (*Request, error) {
	request := &Request{r, time.Now(), nil, nil, nil}
	for _, route := range s.Routes {
		if route.ShouldHandleRequest(r) {
			request.Route = route
		}
	}

	if request.Route != nil {
		var err error
		request.SourceOptions, request.ProcessorOptions, err =
			request.Route.SourceAndProcessorOptionsForRequest(r)
		if err != nil {
			return nil, err
		}
	}

	return request, nil
}

// ResponseWriter is a wrapper around http.ResponseWriter that provides
// access to the response status and size after they have been set.
type ResponseWriter struct {
	w      http.ResponseWriter
	Status int
	Size   int
}

// NewResponseWriter creates a new ResponseWriter by wrapping http.ResponseWriter.
func (s *Server) NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w: w}
}

// WriteHeader forwards to http.ResponseWriter's WriteHeader method.
func (hw *ResponseWriter) WriteHeader(status int) {
	hw.Status = status
	hw.w.WriteHeader(status)
}

// SetHeader sets the value for a response header.
func (hw *ResponseWriter) SetHeader(name, value string) {
	hw.w.Header().Set(name, value)
}

// Writes data the output stream.
func (hw *ResponseWriter) Write(data []byte) (int, error) {
	hw.Size += len(data)
	return hw.w.Write(data)
}

// WriteError writes an error response.
func (hw *ResponseWriter) WriteError(message string, status int) {
	hw.SetHeader("Content-Type", "text/plain; charset=utf-8")
	hw.WriteHeader(status)
	hw.Write([]byte(message))
}

// WriteImage writes an image to the output stream and sets the appropriate headers.
func (hw *ResponseWriter) WriteImage(image *Image) {
	hw.SetHeader("Content-Type", image.MIMEType())
	hw.SetHeader("Content-Length", fmt.Sprintf("%d", image.Size()))
	hw.WriteHeader(http.StatusOK)
	hw.Write(image.Bytes())
}
