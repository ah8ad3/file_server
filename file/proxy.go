package file

import "net/http"

type FileProxy interface {
	Download(http.ResponseWriter, *http.Request) error
}
