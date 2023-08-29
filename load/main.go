// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
)

var (
	source = flag.String("source", "./README.md", "path to markdown file to upload")
	addr   = flag.String("addr", "http://localhost:8080", "address of server")

	count = flag.Int("count", math.MaxInt, "Number of requests to send")
)

// generateLoad sends count requests to the server.
func generateLoad(count int) error {
	if *source == "" {
		return fmt.Errorf("-source must be set to a markdown source file")
	}
	if *addr == "" {
		return fmt.Errorf("-addr must be set to the address of the server (e.g., http://localhost:8080)")
	}

	src, err := os.ReadFile(*source)
	if err != nil {
		return fmt.Errorf("error reading source: %v", err)
	}
	reader := bytes.NewReader(src)

	url := *addr + "/render"

	for i := 0; i < count; i++ {
		reader.Seek(0, io.SeekStart)

		resp, err := http.Post(url, "text/markdown", reader)
		if err != nil {
			return fmt.Errorf("error writing request: %v", err)
		}
		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		resp.Body.Close()
	}

	http.Get(*addr + "/quit")

	return nil
}

func main() {
	flag.Parse()

	if err := generateLoad(*count); err != nil {
		log.Fatal(err)
	}
}
