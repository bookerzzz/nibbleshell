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
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bookerzzz/s3"
)

const (
	ImageSourceTypeS3 ImageSourceType = "s3"
)

type S3ImageSource struct {
	Config *SourceConfig
}

func NewS3ImageSourceWithConfig(config *SourceConfig) (ImageSource, error) {
	return &S3ImageSource{
		Config: config,
	}, nil
}

func (s *S3ImageSource) GetImage(request *ImageSourceOptions) (*Image, error) {
	httpRequest, err := s.signedHTTPRequestForRequest(request)
	if err != nil {
		return nil, err
	}
	httpResponse, err := http.DefaultClient.Do(httpRequest)
	defer httpResponse.Body.Close()
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("Error downlading image (url=%v)", httpRequest.URL)
	}
	image, err := NewImageFromBuffer(httpResponse.Body)
	if err != nil {
		// attempt to consume all body anyway, ignore errors
		_, _ = ioutil.ReadAll(httpResponse.Body)

		// return image read error
		return nil, err
	}
	return image, nil
}

func (s *S3ImageSource) signedHTTPRequestForRequest(request *ImageSourceOptions) (*http.Request, error) {
	path := s.Config.Directory + request.Path
	imageURLPathComponents := strings.Split(path, "/")

	for index, component := range imageURLPathComponents {
		component = url.QueryEscape(component)
		imageURLPathComponents[index] = component
	}
	requestURL := &url.URL{
		Opaque: strings.Join(imageURLPathComponents, "/"),
		Scheme: "http",
		Host:   fmt.Sprintf("%s.s3.amazonaws.com", s.Config.S3Bucket),
	}

	httpRequest, err := http.NewRequest("GET", requestURL.RequestURI(), nil)
	if err != nil {
		return nil, err
	}
	httpRequest.URL = requestURL
	httpRequest.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	s3.Sign(httpRequest, s3.Keys{
		AccessKey: s.Config.S3AccessKey,
		SecretKey: s.Config.S3SecretKey,
	})

	return httpRequest, nil
}

func init() {
	RegisterSource(ImageSourceTypeS3, NewS3ImageSourceWithConfig)
}
