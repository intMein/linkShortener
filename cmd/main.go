package main

import (
	"fmt"
	"io"
	"linkShortener/database"
	"linkShortener/utils"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const shortenURLLength = 5

var storage database.Shortener

type Handler struct{}

func generateUniqueURL() string {
	for {
		u := utils.GenURL(shortenURLLength)

		ok, err := storage.Has(u)

		if err != nil {
			panic("failed to access database")
		}

		if !ok {
			return u
		}
	}
}

func registerHandler(w http.ResponseWriter, req *http.Request) {
	res, err := io.ReadAll(req.Body)

	if err != nil {
		slog.Error("registerHandler: ", err)
		return
	}

	u := string(res)

	parsedUrl, err := url.ParseRequestURI(u)

	if err != nil {
		slog.Error("registerHandler: invalid URL provided!")
		return
	}

	if parsedUrl.RawQuery == "" && parsedUrl.Path == "" && !strings.HasSuffix(u, "/") {
		u = u + "/"
	}

	shorten := generateUniqueURL()

	err = storage.Set(shorten, u)

	if err != nil {
		slog.Error("registerHandler: failed to register URL: ", err)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", shorten)))

	if err != nil {
		slog.Error("registerHandler: failed to write response")
	}
}

func (_ *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" {
		registerHandler(w, req)
		return
	}

	u, err := storage.Get(path[1:])

	if err != nil {
		return
	}

	_, err = w.Write([]byte(u))

	if err != nil {
		slog.Error("mainHandler: failed to write response")
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		var err error = nil

		storage, err = database.NewPostgresDB("user", "weakpassword", "postgres")

		if err != nil {
			log.Fatalln("failed to open Postgres")
		}

	} else {
		storage = database.NewMemoryDB()
	}

	log.Fatalln(http.ListenAndServe("0.0.0.0:8080", &Handler{}))
}
