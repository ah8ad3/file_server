package file

import "net/http"

type FileHandler interface {
	GetFile(http.ResponseWriter, *http.Request)
}
