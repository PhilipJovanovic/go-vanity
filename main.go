package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	htmlTemplate = `<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="{{.URL}} git https://{{.GitHub}}" />
        <meta name="go-source" content="{{.URL}} https://{{.GitHub}} https://{{.GitHub}}/tree/master{/dir} https://{{.GitHub}}/tree/master{/dir}/{file}#L{line}"/>
        <meta http-equiv="refresh" content="0; url=https://{{.Pkg}}" />
    </head>
    <body>
        Nothing to see here. Please <a href="https://{{.Pkg}}">move along</a>.
    </body>
</html>`
	BASE_URL   = os.Getenv("BASE_URL")
	GITHUB_URL = os.Getenv("GITHUB_URL")
)

type PageData struct {
	URL    string
	GitHub string
	Pkg    string
}

func handler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	data := PageData{
		URL:    fmt.Sprintf("%s/%s", BASE_URL, slug),
		Pkg:    fmt.Sprintf("pkg.go.dev/%s/%s", BASE_URL, slug),
		GitHub: fmt.Sprintf("%s/%s", GITHUB_URL, slug),
	}

	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func main() {
	if BASE_URL == "" {
		log.Fatal("BASE_URL environment variable must be set")
	}

	if GITHUB_URL == "" {
		log.Fatal("GITHUB_URL environment variable must be set")
	}

	port := ":1337"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	http.HandleFunc("/", handler)

	fmt.Printf("Start server at http://localhost%s\n", port)
	fmt.Printf("BASE_URL: %s\n", BASE_URL)
	fmt.Printf("GITHUB_URL: %s\n", GITHUB_URL)
	log.Fatal(http.ListenAndServe(port, nil))
}
