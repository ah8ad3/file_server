package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/ah8ad3/file_server/api"
	"github.com/ah8ad3/file_server/file"
	"github.com/ah8ad3/file_server/permission"
	"github.com/ah8ad3/file_server/proxy"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	pemissionUrl := os.Getenv("PERMISSION_URL")
	adminToken := os.Getenv("ADMIN_TOKEN")
	per, err := permission.NewPermission(pemissionUrl, adminToken, 5)
	if err != nil {
		log.Fatal(err)
	}
	proxyPass := os.Getenv("PROXY_PASS")
	openAccessScopes := os.Getenv("OPEN_ACCESS_SCOPES")
	pr := proxy.NewProxy(proxyPass)
	service := file.NewFileService(per, pr)

	openScopes := make([]string, 0)
	splittedPath := strings.Split(openAccessScopes, ",")
	for _, scope := range splittedPath {
		openScopes = append(openScopes, scope)
	}
	handler := api.NewFileHandler(service, openScopes)

	router := mux.NewRouter()
	router.Path("/file/open").Queries("path", "{path}").HandlerFunc(handler.GetOpenAccessFile).Methods("GET")
	router.Path("/file").Queries("path", "{path}").HandlerFunc(handler.GetFile).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	fileServerPort := os.Getenv("FILESERVER_PORT")
	if fileServerPort == "" {
		fileServerPort = "8888"
	}
	srv := &http.Server{
		Addr:         "0.0.0.0:" + fileServerPort,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		Handler:      handlers.CORS(headersOk, originsOk, methodsOk)(router),
	}
	srv.SetKeepAlivesEnabled(false)

	fmt.Println("Listening at :" + fileServerPort)
	log.Fatal(srv.ListenAndServe())

}
