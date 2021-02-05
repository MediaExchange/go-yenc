# go-yenc

[![GoDoc](https://godoc.org/github.com/mediaexchange/nazbaz/github?status.svg)](https://godoc.org/github.com/mediaexchange/go-nntp)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go version](https://img.shields.io/badge/go-~%3E1.15-green.svg)](https://golang.org/doc/devel/release.html#go1.15)

`go-yenc` is a yenc encoder and decoder library for Go.

## Usage

**Encoding**

```go
yenc.Encode(out io.Writer, in io.Reader, lineLength int, filename string) error
```

`Encode` uses a basic `io.Reader` to access the binary data to be encoded and
an `io.Writer` that receives the encoded inforation.

The `lineLength` parameter determines the length of each encoded line. Per
[version 1.3 of the specification](http://www.yenc.org/yEnc-draft-1.txt), the
line length may be one character longer than this as an escape code cannot
appear in isolation.

`fileName` contains the name of the original file being encoded. As the
function utilizes low-level readers and writers, the name must be provided
separately. The file name is encoded into the `=ybegin` header and is emitted
by the `Decode` function

**Decoding**

```go
yenc.Decode(out io.Writer, in io.Reader) error
```

The current implementation of `Decode` currently handles only single-part
files. There are several improvements to the code underway.

Similar to `Encode`, `Decode` accepts a basic `io.Reader` which supplies the
encoded data and a basic `io.Writer` which receives the decoded infomration.

## Contributing

 1.  Fork it
 2.  Create a feature branch (`git checkout -b new-feature`)
 3.  Commit changes (`git commit -am "Added new feature xyz"`)
 4.  Push the branch (`git push origin new-feature`)
 5.  Create a new pull request.

## Maintainers

* [Media Exchange](http://github.com/MediaExchange)

## License

    Copyright 2020 MediaExchange.io
     
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
        http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
