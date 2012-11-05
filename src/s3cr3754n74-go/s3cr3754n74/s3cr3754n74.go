// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package s3cr3754n74

import (
	"io"
	"net/http"
	"text/template"

	"appengine"
)

func serve404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Not Found")
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Internal Server Error")
	c.Errorf("%v", err)
}

var mainPage = template.Must(template.New("secretsanta").Parse(
	`<html><body>
<form action="/sign" method="post">
<div><textarea name="content" rows="30" cols="40"></textarea></div>
<div><input type="submit" value="Send e-mails"></div>
</form></body></html>
`))

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" || r.URL.Path != "/" {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mainPage.Execute(w, nil); err != nil {
		c.Errorf("%v", err)
	}
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		serve404(w)
		return
	}
	c := appengine.NewContext(r)
	if err := r.ParseForm(); err != nil {
		serveError(c, w, err)
		return
	}
	emails := r.FormValue("content")
	http.Redirect(w, r, "/", http.StatusFound)
}

func init() {
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/sendemails", handleSend)
}
