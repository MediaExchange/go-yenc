/*
   Copyright 2021 MediaExchange.io

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package yenc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"strconv"
	"strings"
)

const (
	Answer byte = 0x2A // 42
	Escape byte = 0x3D // 61
	Crit   byte = 0x40 // 64

	Ybegin string = "=ybegin"
	Yend   string = "=yend"
	Ypart  string = "=ypart"
)

// Encode writes ASCII characters that were encoded from binary input using
// the yEnc algorithm, made popular by binary Usenet groups. Any line length
// may be used, but the most commonly used value is 128 characters.
func Encode(out io.Writer, in io.Reader, lineLength int, filename string) error {
	// Use buffering for efficiency.
	rd := bufio.NewReader(in)

	// A CRC32 is computed on the input, but there's no way to reset a reader
	// back to the beginning, so it's copied here as it's processed. This is
	// not the most efficient use of memory.
	var inBuf bytes.Buffer

	// The first line of yEnc output includes the length of the input, which
	// can't be determined from an `io.Reader` in Go. Therefore, the output
	// must be buffered until the length of the input is known.
	var outBuf bytes.Buffer

	var b byte = 0
	var err error = nil

	// Outer loop ends when the Reader has no more bytes (read/EOF error)
	for err == nil {
		// Inner loop is for the length of the emitted line.
		for l := 0; l < lineLength; l++ {
			// Get a byte
			b, err = rd.ReadByte()
			if err != nil {
				break
			}

			inBuf.WriteByte(b)

			if IsCritical(b + Answer) {
				outBuf.WriteByte(Escape)
				b += Crit
			} else {
				b += Answer
			}

			outBuf.WriteByte(b)
		}

		outBuf.WriteByte('\r')
		outBuf.WriteByte('\n')
	}

	// Only inform the caller of errors other than reaching the end of the `io.Reader`
	if err != io.EOF {
		return err
	}

	// Write the encoded data.
	_, err = fmt.Fprintf(out, "%s line=%d size=%d name=%s\r\n", Ybegin, lineLength, inBuf.Len(), filename)
	if err != nil {
		return err
	}

	_, err = outBuf.WriteTo(out)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "%s size=%d crc32=%08x\r\n", Yend, inBuf.Len(), crc32.ChecksumIEEE(inBuf.Bytes()))
	return err
}

// Decode converts ASCII text to binary data using the yEnd algorithm.
func Decode(out io.Writer, in io.Reader) (err error) {
	err = nil

	// Buffer the input and output
	rd := bufio.NewReader(in)
	wr := bufio.NewWriter(out)

	var outLength uint32 = 0
	var isEscape bool = false
	var header map[string]string = nil
	//var part   map[string]string = nil
	var footer map[string]string = nil

	for err == nil {
		buf, err := ReadLine(rd)
		if err != nil {
			break
		}

		line := buf.String()

		// Ignore lines that appear before the ybegin marker line
		if header == nil && !strings.HasPrefix(line, "=ybegin") {
			continue
		}

		// Ignore lines that appear after the yend marker line
		if footer != nil {
			continue
		}

		// Process the header if it hasn't already been encountered.
		if strings.HasPrefix(line, Ybegin) {
			if header != nil {
				err = errors.New("ybegin marker line found multiple times")
				break
			}
			header = ParseMarkerLine(line, Ybegin)
			continue
		}

		// Process the multi-part header
		if strings.HasPrefix(line, Ypart) {
			// Need to handle this...
			//part = ParseMarkerLine(line, Ypart)
			continue
		}

		// Process the footer if it hasn't already been encountered.
		if strings.HasPrefix(line, Yend) {
			if footer != nil {
				err = errors.New("yend marker line found multiple times")
				break
			}
			if header == nil {
				err = errors.New("yend marker line cannot appear before ybegin marker line")
				break
			}
			footer = ParseMarkerLine(line, Yend)
			continue
		}

		// Process encoded lines
		for _, b := range buf.Bytes() {
			if b == Escape {
				isEscape = true
				continue
			}

			if isEscape {
				isEscape = false
				b-=Crit
			} else {
				b -= Answer
			}

			if err = wr.WriteByte(b); err != nil {
				// The outer loop will terminate due to the err != nil condition
				break
			}
			outLength++
		}
	}

	wr.Flush()

	// Validation time
	if header["size"] != footer["size"] {
		err = errors.New(fmt.Sprintf("ybegin size [%s] and yend size [%s] do not match", header["size"], footer["size"]))
		return
	}

	if header["size"] != strconv.FormatUint(uint64(outLength), 10) {
		err = errors.New(fmt.Sprintf("read size [%s] and stored size [%d] do not match", header["size"], outLength))
	}

	return
}

func ReadLine(in *bufio.Reader) (buf bytes.Buffer, err error) {
	var isPrefix bool = true
	for isPrefix && err == nil {
		var line []byte
		line, isPrefix, err = in.ReadLine()
		buf.Write(line)
	}

	return
}

// Parses a marker line that uses the format, `prefix key=value key=value...`
func ParseMarkerLine(line string, prefix string) map[string]string {
	marker := make(map[string]string)
	line = strings.TrimSpace(strings.TrimPrefix(line, prefix))
	entries := strings.Split(line, " ")
	for _, e := range entries {
		parts := strings.Split(e, "=")
		marker[parts[0]] = parts[1]
	}
	return marker
}

// IsCritical returns `true` if the byte is one of the "critical" values.
func IsCritical(b byte) bool {
	return b == 0x00 || b == 0x0A || b == 0x0D || b == Escape
}
