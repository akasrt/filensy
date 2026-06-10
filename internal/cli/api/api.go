package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/akasrt/filensy/internal/config/userconfig"
	"github.com/akasrt/filensy/internal/util/httputil"
)

const (
	host    = "http://127.0.0.1:8080"
	fileURL = host + "/api/v1" + "/file"
)

var client = &http.Client{}

func PostFile(name, ttl, visibility string, isEncrypted bool, file io.Reader) (int, httputil.Response, error) {
	var apiResp httputil.Response

	req, err := http.NewRequest(http.MethodPost, fileURL, file)
	if err != nil {
		return 0, apiResp, err
	}

	setAuthHeader(req)
	req.Header.Set("Content-Type", "application/octet-stream")

	q := req.URL.Query()
	q.Add("name", name)
	q.Add("enc", strconv.FormatBool(isEncrypted))

	if ttl != "" {
		q.Add("ttl", ttl)
	}
	if visibility != "" {
		q.Add("visibility", visibility)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return 0, apiResp, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return resp.StatusCode, apiResp, err
	}

	return resp.StatusCode, apiResp, nil
}

func GetFile(code, token string) (int, httputil.Response, io.ReadCloser, map[string]string, error) {
	url := fmt.Sprintf("%s/%s", fileURL, code)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, httputil.Response{}, nil, nil, err
	}
	setAuthHeader(req)
	if token != "" {
		req.Header.Set("X-File-Token", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, httputil.Response{}, nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		var apiResp httputil.Response
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return resp.StatusCode, httputil.Response{}, nil, nil, err
		}
		return resp.StatusCode, apiResp, nil, nil, nil
	}

	headers := map[string]string{
		"X-File-Name":         resp.Header.Get("X-File-Name"),
		"X-Is-Encrypted":      resp.Header.Get("X-Is-Encrypted"),
		"Content-Disposition": resp.Header.Get("Content-Disposition"),
		"Content-Length":      resp.Header.Get("Content-Length"),
	}

	return resp.StatusCode, httputil.Response{}, resp.Body, headers, nil
}

func GetFileMetadata(code, token string) (int, httputil.Response, error) {
	var apiResp httputil.Response

	url := fmt.Sprintf("%s/%s/meta", fileURL, code)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, apiResp, err
	}
	setAuthHeader(req)
	if token != "" {
		req.Header.Set("X-File-Token", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, apiResp, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return resp.StatusCode, apiResp, err
	}

	return resp.StatusCode, apiResp, nil
}

func DeleteFile(code string, token string) (int, httputil.Response, error) {
	url := fmt.Sprintf("%s/%s", fileURL, code)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return 0, httputil.Response{}, err
	}
	setAuthHeader(req)
	req.Header.Set("X-File-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return 0, httputil.Response{}, err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
		var apiResp httputil.Response

		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return resp.StatusCode, httputil.Response{}, err
		}

		return resp.StatusCode, apiResp, nil
	}

	return resp.StatusCode, httputil.Response{}, nil
}

func setAuthHeader(req *http.Request) {
	conf := userconfig.GetConfig()
	if conf.AuthKey != "" {
		authToken := fmt.Sprintf("Bearer %s", conf.AuthKey)
		req.Header.Set("Authorization", authToken)
	}
}
