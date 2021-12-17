package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ah8ad3/file_server/file"
)

type proxy struct {
	baseUrl string
}

func NewProxy(baseUrl string) file.FileProxy {
	return &proxy{
		baseUrl: baseUrl,
	}
}

func (p *proxy) Download(w http.ResponseWriter, r *http.Request) error {
	path := r.FormValue("path")
	remote, err := url.Parse(p.baseUrl)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	r.Host = remote.Host

	// minio check this header so we have to remove it
	r.Header.Del("Authorization")

	newRequestUrl, err := url.Parse(path)
	if err != nil {
		return err
	}
	r.URL = newRequestUrl
	proxy.ServeHTTP(w, r)
	return nil
}
