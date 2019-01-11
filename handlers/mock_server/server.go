package main

import (
	"context"
	"log"
	"net/http"

	"github.com/golang/mock/gomock"
)

const addr = "127.0.0.1:8080"

func main() {
	ctx := context.Background()
	ctrl := gomock.NewController(&testReporter{})
	defer ctrl.Finish()

	router := newHandlers(ctx, ctrl).NewRouter()
	http.Handle("/", router)

	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", static)

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
