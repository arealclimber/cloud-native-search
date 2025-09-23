// pprof.go
//go:build debug

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

// StartPprof 在 :6060 開 profiling server
func StartPprof() {
	go func() {
		log.Println("pprof listening on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
