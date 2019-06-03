package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "HTTP port")
	flag.Parse()
}

func main() {
	fmt.Printf("gitpanda started: port=%d\n", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("It works"))
		return
	}
}
