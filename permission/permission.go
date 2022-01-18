package permission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ah8ad3/file_server/file"
)

type permission struct {
	client           http.Client
	timeout          int
	permissionServer string
	adminToken       string
}

func newPermission(permissonServer string) (http.Client, error) {
	client := &http.Client{}
	return *client, nil
}

func NewPermission(permissonServer, adminToken string, timeout int) (file.FilePermission, error) {
	client, err := newPermission(permissonServer)
	if err != nil {
		return nil, err
	}
	return &permission{
		client:           client,
		timeout:          timeout,
		permissionServer: permissonServer,
		adminToken:       adminToken,
	}, nil
}

func (p *permission) HasPermission(f file.File, jwt string) (bool, error) {
	// ignore permission check if admin requests
	if jwt == p.adminToken {
		return true, nil
	}
	jsonValue, err := json.Marshal(f)
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest("POST", p.permissionServer, bytes.NewBuffer(jsonValue))
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", jwt))
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 401:
		return false, file.ErrJwtNotProvided
	case 403:
		return false, file.ErrFilePermissionDenied
	case 404:
		return false, file.ErrFileNotFound
	case 200:
		return true, nil
	}

	return false, nil
}
