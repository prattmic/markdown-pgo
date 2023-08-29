// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"

	"gitlab.com/golang-commonmark/markdown"
)

func render(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	src, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	md := markdown.New(
		markdown.XHTMLOutput(true),
		markdown.Typographer(true),
		markdown.Linkify(true),
		markdown.Tables(true),
	)

	var buf bytes.Buffer
	if err := md.Render(&buf, src); err != nil {
		log.Printf("error converting markdown: %v", err)
		http.Error(w, "Malformed markdown", http.StatusBadRequest)
		return
	}

	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("error writing response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	f, err := os.Create("/tmp/cpu.pprof")
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	http.HandleFunc("/quit", func(http.ResponseWriter, *http.Request) {
		pprof.StopCPUProfile()
		f.Close()

		f, err := os.Create("/tmp/heap.pprof")
		if err != nil {
			panic(err)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}
		f.Close()

		os.Exit(0)
	})

	http.HandleFunc("/render", render)
	log.Printf("Serving on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

