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

package main

import (
	"flag"
	"fmt"
	"github.com/MediaExchange/go-yenc"
	"os"
	"path/filepath"
	"strings"
)

var (
	Output string
)

func main() {
	flag.StringVar(&Output, "o", "", "Full pathname to output file")
	flag.Parse()

	if len(os.Args) < 3 {
		help()
		return
	}

	switch strings.ToLower(flag.Arg(0)) {
	case "encode":
		encode()
	case "decode":
		decode()
	case "help":
		help()
	default:
		help()
	}
}

func encode() {
	filename := flag.Arg(1)
	if !fileExists(filename) {
		fmt.Printf("file [%s] does not exist", filename)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing input file: %s", err)
		}
	}()

	yencName := Output
	if len(yencName) == 0 {
		yencName = filename + ".yenc"
	}

	if fileExists(yencName) {
		fmt.Printf("file [%s] exists; cannot overwrite\n", yencName)
		return
	}

	yencFile, err := os.Create(yencName)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	defer func() {
		err := yencFile.Close()
		if err != nil {
			fmt.Printf("error closing output file: %s", err)
		}
	}()

	yenc.Encode(yencFile, file, 128, filepath.Base(filename))
}

func decode() {
	filename := os.Args[2]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	yenc.Decode(os.Stdout, file)
}

// help displays the program's sub-commands and arguments.
func help() {
	fmt.Println("Usage:")
	fmt.Print("\n")
	fmt.Println("    yenc <command> <filename>")
	fmt.Print("\n")
	fmt.Println("The commands are:")
	fmt.Print("\n")
	fmt.Println("    encode	Encode a file")
	fmt.Println("    decode	Decode a file")
	fmt.Print("\n")
	flag.PrintDefaults()
}

func fileExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}