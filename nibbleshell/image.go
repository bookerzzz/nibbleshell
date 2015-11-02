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
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

type imageFormat string

const (
	JPEG imageFormat = "jpeg"
	PNG  imageFormat = "png"
)

func (f imageFormat) String() string {
	return "image/" + string(f)
}

type Image struct {
	image  image.Image
	format imageFormat
	buffer bytes.Buffer
}

func (i *Image) MIMEType() string {
	return i.format.String()
}

func (i *Image) Size() int {
	return i.buffer.Len()
}

func (i *Image) Bytes() []byte {
	return i.buffer.Bytes()
}

func NewImageFromFile(file *os.File) (*Image, error) {
	return NewImageFromBuffer(file)
}

func NewImageFromBuffer(i io.ReadCloser) (*Image, error) {
	var format string
	var err error
	img := &Image{}
	img.image, format, err = image.Decode(i)
	if err != nil {
		return nil, err
	}
	switch format {
	case "jpeg":
		img.format = JPEG
	case "png":
		img.format = PNG
	default:
		return nil, errors.New("image format not supported")
	}

	return img, nil
}
