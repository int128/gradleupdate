package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang/mock/gomock"
)

const addr = "127.0.0.1:8080"

func main() {
	ctx := context.Background()
	ctrl := gomock.NewController(&testReporter{})
	defer ctrl.Finish()
	router := newHandlers(ctx, ctrl).NewRouter()
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/":
			http.ServeFile(w, r, "static/index.html")
		case strings.HasPrefix(r.URL.Path, "/static"):
			static.ServeHTTP(w, r)
		default:
			router.ServeHTTP(w, r)
		}
	})

	log.Printf("Open http://%s", addr)
	if err := http.ListenAndServe(addr, http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}

type testReporter struct{}

func (t *testReporter) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (t *testReporter) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
