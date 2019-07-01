package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

var host = "s2-us2.startpage.com"

var ownhost = "mirror.loerfy.now.sh"

// var ownhost = "127.0.0.1:3000"

func hasGziped(coding string) bool {
	return strings.HasPrefix(coding, "gz")
}

func isTextType(typeName string) bool {
	return strings.HasPrefix(typeName, "text") || strings.HasPrefix(typeName, "app")
}

func rewriteBody(resp *http.Response) (err error) {
	cType := resp.Header.Get("Content-Type")
	cEncoding := resp.Header.Get("Content-Encoding")
	StatusCode := resp.StatusCode
	var b []byte
	if hasGziped(cEncoding) {
		reader, _ := gzip.NewReader(resp.Body)
		b, err = ioutil.ReadAll(reader)
	} else {
		b, err = ioutil.ReadAll(resp.Body)
	}
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	if hasGziped(cEncoding) {
		resp.Header.Del("Content-Encoding")
	}
	if isTextType(cType) {
		b = bytes.Replace(b, []byte("www.startpage.com"), []byte(host), -1) // replace html
		b = bytes.Replace(b, []byte(host), []byte(ownhost), -1)             // replace html
	}
	if StatusCode == 302 || StatusCode == 301 {
		lo := resp.Header.Get("location")
		newLo := strings.ReplaceAll(lo, "www.startpage.com", ownhost)
		resp.Header.Set("Location", newLo)
		cookie := strings.ReplaceAll(resp.Header.Get("set-cookie"), "domain=startpage.com;", "")
		resp.Header.Set("Set-Cookie", cookie)
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	rpURL, err := url.Parse("https://" + host)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(rpURL)
	proxy.ModifyResponse = rewriteBody
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		r.Host = host
	}
	proxy.ServeHTTP(w, r)
}

// func main() {
// 	http.HandleFunc("/", Handler)
// 	// http.ListenAndServe(":3000", nil)
// 	http.ListenAndServeTLS(":3000", "/home/wincer/.local/share/mkcert/rootCA.pem", "/home/wincer/.local/share/mkcert/rootCA-key.pem", nil)
// 	log.Println("Listening in :3000")
// }
