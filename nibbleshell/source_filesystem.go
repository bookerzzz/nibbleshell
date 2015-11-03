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
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	ImageSourceTypeFilesystem ImageSourceType = "filesystem"
)

type FileSystemImageSource struct {
	Config *SourceConfig
}

func NewFileSystemImageSourceWithConfig(config *SourceConfig) (ImageSource, error) {
	source := &FileSystemImageSource{
		Config: config,
	}

	baseDirectory, err := os.Open(source.Config.Directory)
	if os.IsNotExist(err) {
		return source, errors.New("source directory does not exist")
	}
	if err != nil {
		return source, err
	}

	fileInfo, err := baseDirectory.Stat()
	if err != nil || !fileInfo.IsDir() {
		return source, errors.New("source directory is not a directory")
	}

	return source, nil
}

func (s *FileSystemImageSource) GetImage(request *ImageSourceOptions) (*Image, error) {
	fileName := s.fileNameForRequest(request)

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	image, err := NewImageFromFile(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *FileSystemImageSource) fileNameForRequest(request *ImageSourceOptions) string {
	// Remove the leading / from the file name and replace the
	// directory separator (/) with something safe for file names (_)
	return filepath.Join(s.Config.Directory, strings.Replace(strings.TrimLeft(request.Path, string(filepath.Separator)), string(filepath.Separator), "_", -1))
}

func init() {
	RegisterSource(ImageSourceTypeFilesystem, NewFileSystemImageSourceWithConfig)
}
