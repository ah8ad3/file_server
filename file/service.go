package file

import "net/http"

type FileService interface {
	HasPermission(File, string) (bool, error)
	ReverceProxy(http.ResponseWriter, *http.Request) error
}
