// Copyright 2015 The quotesrv Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "./statik/"
	"bytes"
	"flag"
	"fmt"
	"github.com/jroimartin/orujo"
	"github.com/jroimartin/orujo-handlers/basic"
	olog "github.com/jroimartin/orujo-handlers/log"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"text/template"
)

var (
	addr       = flag.String("addr", ":7000", "HTTP service address")
	quotesFile = flag.String("quotesfile", "quotes.txt", "quotes file")
	auth       = flag.Bool("auth", false, "enable basic authentication")
	user       = flag.String("user", "user", "basic auth username")
	pass       = flag.String("pass", "s3cr3t", "basic auth password")
	tls        = flag.Bool("tls", false, "enable TLS")
	certFile   = flag.String("cert", "cert.pem", "certificate file")
	keyFile    = flag.String("key", "key.pem", "private key file")
	re         = regexp.MustCompile(`[\r\n]+`)
)

var mutex sync.RWMutex

const errorMessage = `
`

func parseTemplate(s string) string {
	type Inventory struct {
		Host string
		Port uint
	}
	data := Inventory{"localhost", 7000}
	tmpl, err := template.New("test").Parse(s)
	if err != nil {
		panic(err)
	}
	msg := &bytes.Buffer{}
	if err := tmpl.Execute(msg, data); err != nil {
		return ""
	}
	return msg.String()
}

var authHandler basic.BasicHandler
var logHandler olog.LogHandler

func registerGetRoute(s *orujo.Server, path string, hf http.HandlerFunc) {
	prefix := path
	if path != "" {
		//prefix = "^/" + path + "/.*"
		prefix = "^/" + path + ".*"
	} else {
		prefix = "^/"
	}
	s.Route(prefix,
		authHandler,
		http.HandlerFunc(hf),
		orujo.M(logHandler),
	).Methods("GET")
}

var statikFS http.FileSystem

func main() {
	flag.Parse()
	newfs, err := fs.New()
	if err != nil {
		fmt.Println("fuck id", err)
	}
	statikFS = newfs
	fmt.Println(statikFS)
	*auth = true
	*user = "a"
	*pass = "b"
	s := orujo.NewServer(*addr)

	logger := log.New(os.Stdout, "[lws] ", log.LstdFlags)
	logHandler = olog.NewLogHandler(logger, logLine)

	authHandler = basic.NewBasicHandler("LocalWebService", *user, *pass)
	authHandler.ErrorMsg = func(w http.ResponseWriter, provuser string) {
		foo, _ := getFile("/error.html")
		fmt.Fprintln(w, parseTemplate(foo))
	}
	registerGetRoute(s, "cmd", handleCmd)
	registerGetRoute(s, "dir", handleDir)
	registerGetRoute(s, "kill", handleKill)
	//	registerGetRoute(s, "start", handleStart)
	//	registerGetRoute(s, "stop", handleStop)
	registerGetRoute(s, "", showIndex)

	s.RouteDefault(http.HandlerFunc(handleGet),
		orujo.M(logHandler))

	logger.Fatalln(s.ListenAndServe())
}

func handleStart(w http.ResponseWriter, r *http.Request) {
}

func handleStop(w http.ResponseWriter, r *http.Request) {
}

func handleList(w http.ResponseWriter, r *http.Request) {
}

func handleCmd(w http.ResponseWriter, r *http.Request) {
	const SHELL = "/system/bin/sh"
	path, err := url.QueryUnescape(r.URL.Path)
	if err != nil {
		path = r.URL.Path
	}
	path = path[5:] // assume "/cmd/"
	if path == "" {
		str, _ := getFile("/cmd.html")
		fmt.Fprintln(w, parseTemplate(str))
	} else {
		fmt.Println(path)
		//argv := strings.Split(path, " ")
		out, err := exec.Command(SHELL, "-c", path).Output()
		// out, err := exec.Command(argv[0], argv[1:]...).Output()
		if err == nil {
			fmt.Fprintln(w, string(out))
		}
	}
}

func handleKill(w http.ResponseWriter, r *http.Request) {
	const SHELL = "/system/bin/sh"
	const KILL = "/system/bin/kill"
	path, err := url.QueryUnescape(r.URL.Path)
	if err != nil {
		path = r.URL.Path
	}
	path = path[6:] // assume "/cmd/"
	if path == "" {
		str, _ := getFile("/cmd.html")
		fmt.Fprintln(w, parseTemplate(str))
	} else {
		fmt.Println(path)
		if path == "b2g" {
			out, err := exec.Command(SHELL, "-c",
				"ps| grep b2g|grep root| grep ' 1 '").Output()
			// out, err := exec.Command(argv[0], argv[1:]...).Output()
			if err == nil {
				outstr := strings.TrimSpace(string(out))
				words := strings.Fields(outstr)
				fmt.Fprintln(w, outstr)
				fmt.Println("kill -9 " + words[1])
				_, err2 := exec.Command(SHELL, "-c", "kill -9 "+words[1]).Output()
				if err2 != nil {
					fmt.Println("ERROR HAPPEN")
				}

			}
		} else {
			fmt.Fprintln(w, "TODO")
		}
	}
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func handleDir(w http.ResponseWriter, r *http.Request) {
	path, err := url.QueryUnescape(r.URL.Path)
	if err != nil {
		path = r.URL.Path
	}
	if len(path) > 4 && path[4] == '/' {
		path = path[4:]
		files, _ := ioutil.ReadDir(path)
		fmt.Fprintln(w, ".", path)
		lastslash := strings.LastIndex(path, "/")
		if lastslash != -1 {
			super_path := path[0:lastslash]
			if super_path == "" {
				super_path = "/"
			}
			fmt.Fprintln(w, ".. <a href='/dir/"+super_path+"'>"+
				super_path+"</a>")
		}
		for _, f := range files {
			fname := f.Name()
			if IsDirectory(path + "/" + fname) {
				fmt.Fprintln(w, "<a href='/dir/"+path+"/"+fname+"'>"+fname+"</a>")
			} else {
				fmt.Fprintln(w, fname)
			}
		}
	} else {
		str, _ := getFile("/dir.html")
		fmt.Fprintln(w, parseTemplate(str))
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	path, err := url.QueryUnescape(r.URL.Path)
	if err != nil {
		path = r.URL.Path
	}
	fmt.Fprintln(w, " --> "+path)
}

/* Statik Filesystem */
func getFile(path string) (string, error) {
	file, err := statikFS.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println("Error")
		return "", err
	}
	res, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func showIndex(w http.ResponseWriter, r *http.Request) {
	str, err := getFile("/index.html")
	if err != nil {
		fmt.Fprintln(w, "404")
	} else {
		fmt.Fprintln(w, str)
	}
}

func errorResponse(w http.ResponseWriter, err error) {
	orujo.RegisterError(w, err)
	w.WriteHeader(http.StatusInternalServerError)
}

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`
