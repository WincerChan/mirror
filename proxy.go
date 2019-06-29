package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

func rewriteBody(resp *http.Response) (err error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	b = bytes.Replace(b, []byte("jxnu.edu"), []byte("xxx"), -1) // replace html
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func main() {
	rpURL, err := url.Parse(os.Args[1])
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(rpURL)
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		r.Host = "www.baidu.com"
	}
	proxy.ModifyResponse = rewriteBody

	http.Handle("/", proxy)
	http.ListenAndServe(":5600", nil)
}
