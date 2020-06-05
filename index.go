package mirror

import (
	"bytes"
	"compress/gzip"
	"github.com/dsnet/compress/brotli"
	"io/ioutil"
	C "mirror/config"
	T "mirror/tool"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var initial bool

func replaceText(text []byte) []byte {
	for _, url := range C.Config.ReplacedURLs {
		text = bytes.ReplaceAll(text,
			[]byte(url.Old), []byte(url.New))
	}
	return text
}

func addNotification(text []byte) []byte {
	if bytes.Contains(text, []byte("doctype")) {
		var script = `<script>
		let SIvCob = document.querySelector('#SIvCob');
		if (SIvCob) SIvCob = document.querySelector('#SIvCob').innerHTML = '这是一个 Google 的镜像站，原理<a target=\'_blank\' href=\'https://blog.itswincer.com/posts/1352252a/\'>戳我</a>'
		</script></body></html>`
		text = bytes.ReplaceAll(text, []byte("</body></html>"), []byte(script))
	}
	return text
}

func replaceRedirect(header http.Header) string {
	domain := regexp.MustCompile(`://(.*?)/`)
	location := header.Get("Location")
	host := domain.FindStringSubmatch(location)[1]
	if host != C.Config.Host.Self {
		return strings.ReplaceAll(location, host, C.Config.Host.Self)
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
	cType := resp.Header.Get("content-type")
	cEncoding := resp.Header.Get("content-encoding")
	StatusCode := resp.StatusCode
	cookie := resp.Header.Get("set-cookie")

	if T.HasGziped(cEncoding) {
		resp.Header.Del("content-encoding")
		reader, _ := gzip.NewReader(resp.Body)
		content, err = ioutil.ReadAll(reader)
	} else if T.HasBrotli(cEncoding) {
		resp.Header.Del("content-encoding")
		reader, _ := brotli.NewReader(resp.Body, &brotli.ReaderConfig{})
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

	content = addNotification(content)
	resp.Body = ioutil.NopCloser(bytes.NewReader(content))
	resp.ContentLength = int64(len(content))
	resp.Header.Set("Content-Length", strconv.Itoa(len(content)))

	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if !initial {
		C.LoadConfig()
		initial = true
	}
	rpURL, err := url.Parse(C.Protocal + C.Config.Host.Proxy)
	T.CheckErr(err)
	proxy := httputil.NewSingleHostReverseProxy(rpURL)
	proxy.ModifyResponse = rewriteBody
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		r.Host = C.Config.Host.Proxy
	}
	proxy.ServeHTTP(w, r)
}

// func main() {
// 	http.HandleFunc("/", Handler)
// 	// http.ListenAndServe(":3000", nil)
// 	http.ListenAndServeTLS(":3000", "/home/wincer/.local/share/mkcert/rootCA.pem", "/home/wincer/.local/share/mkcert/rootCA-key.pem", nil)
// 	log.Println("Listening in :3000")
// }
