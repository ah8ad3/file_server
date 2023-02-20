package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ah8ad3/file_server/file"
)

type fileHandler struct {
	fileService   file.FileService
	allowedScopes []string
}

func NewFileHandler(fileService file.FileService, allowedScopes []string) file.FileHandler {
	return &fileHandler{
		fileService:   fileService,
		allowedScopes: allowedScopes,
	}
}

func (f *fileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	authorization = strings.Replace(authorization, "Token ", "", 1)
	path := r.FormValue("path")
	splittedPath := strings.Split(path, "/")
	if len(splittedPath) < 2 {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"status": "%v"}`, "path is not in correct format")))
		return
	}

	ff := file.File{FileName: path, Scope: splittedPath[0]}
	_, err := f.fileService.HasPermission(ff, authorization)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case file.ErrFileNotFound:
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"status": "%v"}`, err.Error())))
			return
		case file.ErrFilePermissionDenied:
			w.WriteHeader(403)
			w.Write([]byte(fmt.Sprintf(`{"status": "%v"}`, err.Error())))
			return
		case file.ErrJwtNotProvided:
			w.WriteHeader(401)
			w.Write([]byte(fmt.Sprintf(`{"status": "%v"}`, err.Error())))
			return
		default:
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`{"status": "%v"}`, err.Error())))
			return
		}
	}
	err = f.fileService.ReverceProxy(w, r)
	if err != nil {
		log.Fatal(err)
	}

}

func (f *fileHandler) GetOpenAccessFile(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	splittedPath := strings.Split(path, "/")
	if len(splittedPath) < 2 {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"status": "%v"}`, "path is not in correct format")))
		return
	}
	requestScope := splittedPath[0]
	for _, scope := range f.allowedScopes {
		if requestScope == scope {
			err := f.fileService.ReverceProxy(w, r)
			if err != nil {
				log.Fatal(err)
			}		
		}
	}
	w.WriteHeader(403)
	w.Write([]byte(fmt.Sprintf(`{"status": "%v"}`, "Scope is not allowed")))
	return
}
