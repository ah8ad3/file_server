package file

import (
	"errors"
	"net/http"
)

var (
	ErrFileNotFound         = errors.New("file not found")
	ErrFilePermissionDenied = errors.New("you can not access this file")
	ErrJwtNotProvided       = errors.New("jwt not found")
)

type fileService struct {
	filePermission FilePermission
	fileProxy      FileProxy
}

func NewFileService(filePermission FilePermission, fileProxy FileProxy) FileService {
	return &fileService{
		filePermission: filePermission,
		fileProxy:      fileProxy,
	}
}

func (f *fileService) HasPermission(fi File, jwt string) (bool, error) {
	return f.filePermission.HasPermission(fi, jwt)
}

func (f *fileService) ReverceProxy(w http.ResponseWriter, r *http.Request) error {
	return f.fileProxy.Download(w, r)
}
