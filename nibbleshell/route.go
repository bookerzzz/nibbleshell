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
	"net/http"
	"regexp"
	"strconv"
)

// A Route handles the internal logic of a Halfshell request. When a request is
// serviced, the appropriate route is chosen after which the image is retrieved
// from the source and processed by the processor.
type Route struct {
	Name           string
	Pattern        *regexp.Regexp
	ImagePathIndex int
	Source         ImageSource
	CacheControl   string
	Statter        Statter
}

// NewRouteWithConfig returns a pointer to a new Route instance created using
// the provided configuration settings.
func NewRouteWithConfig(config *RouteConfig, statterConfig *StatterConfig) (*Route, error) {
	source, err := NewImageSourceWithConfig(config.SourceConfig)
	if err != nil {
		return nil, err
	}
	statter, err := NewStatterWithConfig(config, statterConfig)
	if err != nil {
		return nil, err
	}

	return &Route{
		Name:           config.Name,
		Pattern:        config.Pattern,
		ImagePathIndex: config.ImagePathIndex,
		CacheControl:   config.CacheControl,
		Source:         source,
		Statter:        statter,
	}, nil
}

// ShouldHandleRequest accepts an HTTP request and returns a bool indicating
// whether the route should handle the request.
func (p *Route) ShouldHandleRequest(r *http.Request) bool {
	return p.Pattern.MatchString(r.URL.Path)
}

// SourceAndProcessorOptionsForRequest parses the source and processor options
// from the request.
func (p *Route) SourceAndProcessorOptionsForRequest(r *http.Request) (
	*ImageSourceOptions, *ImageProcessorOptions, error) {

	matches := p.Pattern.FindAllStringSubmatch(r.URL.Path, -1)[0]
	path := matches[p.ImagePathIndex]

	var width, height, x, y uint64
	var err error
	width, err = strconv.ParseUint(r.FormValue("w"), 10, 32)
	if err != nil {
		return nil, nil, err
	}
	height, err = strconv.ParseUint(r.FormValue("h"), 10, 32)
	if err != nil {
		return nil, nil, err
	}

	x, err = strconv.ParseUint(r.FormValue("x"), 10, 32)
	if err != nil {
		return nil, nil, err
	}
	y, err = strconv.ParseUint(r.FormValue("y"), 10, 32)
	if err != nil {
		return nil, nil, err
	}

	var scale_x, scale_y int64
	scale_x, err = strconv.ParseInt(r.FormValue("scale_x"), 10, 32)
	if err != nil {
		return nil, nil, err
	}
	scale_y, err = strconv.ParseInt(r.FormValue("scale_y"), 10, 32)
	if err != nil {
		return nil, nil, err
	}

	if scale_x != 1 && scale_x != -1 {
		return nil, nil, fmt.Errorf("only horizontal flip supported for X scaling")
	}

	if scale_y != 1 {
		return nil, nil, fmt.Errorf("Y scaling not supported")
	}

	return &ImageSourceOptions{Path: path}, &ImageProcessorOptions{
		Width:  uint32(width),
		Height: uint32(height),
		X:      uint32(x),
		Y:      uint32(y),
		ScaleX: int32(scale_x),
		ScaleY: int32(scale_y),
	}, nil
}
