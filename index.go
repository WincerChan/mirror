package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	C "mirror/config"
	"mirror/middleware"
	T "mirror/tool"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func replaceText(text []byte) []byte {
	for _, url := range C.GetConfig().ReplacedURLs {
		text = bytes.ReplaceAll(text,
			[]byte(url.Old), []byte(url.New))
	}
	return text
}

func replaceRedirect(header http.Header) string {
	domain := regexp.MustCompile(`://(.*?)/`)
	location := header.Get("Location")
	host := domain.FindStringSubmatch(location)[1]
	if host != C.GetConfig().Host.Self {
		return strings.ReplaceAll(location, host, C.GetConfig().Host.Self)
	}
	return location
}

func removeCookie(cookie string) string {
	domain := regexp.MustCompile(`(domain=.*?;)`)
	newCookit := domain.ReplaceAllLiteralString(cookie, "")
	return newCookit
}

func rewriteBody(resp *http.Response) (err error) {
	if nil != resp {
		defer resp.Body.Close()
	}
	T.CheckErr(err)

	var content []byte
	cType := resp.Header.Get("Content-Type")
	cEncoding := resp.Header.Get("Content-Encoding")
	StatusCode := resp.StatusCode
	cookie := resp.Header.Get("Set-Cookie")

	if T.HasGziped(cEncoding) {
		resp.Header.Del("Content-Encoding")
		reader, _ := gzip.NewReader(resp.Body)
		content, err = ioutil.ReadAll(reader)
	} else {
		content, err = ioutil.ReadAll(resp.Body)
	}
	T.CheckErr(err)
	if T.IsTextType(cType) {
		content = replaceText(content)
	}
	if StatusCode == 302 || StatusCode == 301 {
		resp.Header.Set("Location", replaceRedirect(resp.Header))
	}
	if cookie != "" {
		resp.Header.Set("Set-Cookie", removeCookie(cookie))
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(content))
	resp.ContentLength = int64(len(content))
	resp.Header.Set("Content-Length", strconv.Itoa(len(content)))

	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	rpURL, err := url.Parse(C.GetConfig().Protocol + C.GetConfig().Host.Proxy)
	T.CheckErr(err)
	proxy := httputil.NewSingleHostReverseProxy(rpURL)
	proxy.ModifyResponse = rewriteBody
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		r.Host = C.GetConfig().Host.Proxy
	}
	proxy.ServeHTTP(w, r)
}

func main() {
	e := &middleware.Engine{}
	e.Use(middleware.RecoverMW())
	e.Use(middleware.IdentityMW())
	e.Use(middleware.CreateHandler(Handler))
	http.HandleFunc("/", e.Run())
	log.Println("Listening in :3000")
	http.ListenAndServe(":3000", nil)
}
