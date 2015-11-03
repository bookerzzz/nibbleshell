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

const StartupTemplateString = `
  _   _ _ _     _     _           _          _ _ 
 | \ | (_) |__ | |__ | | ___  ___| |__   ___| | |
 |  \| | | '_ \| '_ \| |/ _ \/ __| '_ \ / _ \ | |
 | |\  | | |_) | |_) | |  __/\__ \ | | |  __/ | |
 |_| \_|_|_.__/|_.__/|_|\___||___/_| |_|\___|_|_|

Running on process {{.Pid}}

Server settings:
  Port: {{.Config.ServerConfig.Port}}
  Read Timeout: {{.Config.ServerConfig.ReadTimeout}}
  Write Timeout: {{.Config.ServerConfig.WriteTimeout}}

StatsD settings:
  Host: {{.Config.StatterConfig.Host}}
  Port: {{.Config.StatterConfig.Port}}
  Enabled: {{.Config.StatterConfig.Enabled}}

Routes:
{{ range $index, $route := .Routes }}  {{ $route.Name }}:
    Pattern: {{ $route.Pattern }}
{{ end }}
`
