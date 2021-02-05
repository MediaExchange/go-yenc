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
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncode1(t *testing.T) {
	// receives the encoded output
	var actual bytes.Buffer

	// input to encode
	file, err := os.Open("testdata/test1.txt")
	if err != nil {
		t.Error(err)
	}

	// encode it
	err = Encode(&actual, file, 128, "test1.txt")
	if err != nil {
		t.Error(err)
	}

	// Compare to the original
	expected, err := ioutil.ReadFile("testdata/test1.yenc")
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(actual.Bytes(), expected) != 0 {
		t.Fatal("actual and expected are different")
	}
}

func TestDecode1(t *testing.T) {
	// receives the decoded output
	var actual bytes.Buffer

	// input to decode
	file, err := os.Open("testdata/test1.yenc")
	if err != nil {
		t.Error(err)
	}

	// decode it
	err = Decode(&actual, file)
	if err != nil {
		t.Error(err)
	}

	// compare to the original
	expected, err := ioutil.ReadFile("testdata/test1.txt")
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(actual.Bytes(), expected) != 0 {
		t.Fatal("actual and expected are different")
	}
}

func TestEncode2(t *testing.T) {
	// receives the encoded output
	var actual bytes.Buffer

	// input to encode
	file, err := os.Open("testdata/test2.bin")
	if err != nil {
		t.Error(err)
	}

	// encode it
	err = Encode(&actual, file, 128, "test2.bin")
	if err != nil {
		t.Error(err)
	}

	// Compare to the original
	expected, err := ioutil.ReadFile("testdata/test2.yenc")
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(actual.Bytes(), expected) != 0 {
		t.Fatal("actual and expected are different")
	}
}

func TestDecode2(t *testing.T) {
	// receives the decoded output
	var actual bytes.Buffer

	// input to decode
	file, err := os.Open("testdata/test2.yenc")
	if err != nil {
		t.Error(err)
	}

	// decode it
	err = Decode(&actual, file)
	if err != nil {
		t.Error(err)
	}

	// compare to the original
	expected, err := ioutil.ReadFile("testdata/test2.bin")
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(actual.Bytes(), expected) != 0 {
		t.Fatal("actual and expected are different")
	}
}
