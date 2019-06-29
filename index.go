package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func rewriteBody(resp *http.Response) (err error) {
	reader, _ := gzip.NewReader(resp.Body)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	cType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(cType, "text") {
		resp.Header.Del("Content-Encoding")
		fmt.Println(true)
		b = bytes.Replace(b, []byte("test"), []byte("d"), -1) // replace html
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func main() {
	rpURL, err := url.Parse("https://" + os.Args[1])
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(rpURL)
	proxy.ModifyResponse = rewriteBody
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		// r.Host = "www.startpage.com"
		r.Host = os.Args[1]
	}

	http.Handle("/", proxy)
	http.ListenAndServe(":3000", nil)
}
