package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type ApiClient struct {
	baseUrl string
}

func NewApiClient() ApiClient {
	return ApiClient{
		baseUrl: "http://localhost:80/api/v1",
	}
}

func (api *ApiClient) Post(path string, data map[string]interface{}) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(body)
	url := api.baseUrl + path

	fmt.Println("POST", url, payload)

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return nil, err
	}

	fmt.Println("RESPONSE", resp.Status)

	return resp, nil
}

func (api *ApiClient) Get(path string) (*http.Response, error) {
	url := api.baseUrl + path

	fmt.Println("GET", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	fmt.Println("RESPONSE", resp.Status)

	return resp, nil
}

func (api *ApiClient) Put(path string, data map[string]interface{}) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(body)
	url := api.baseUrl + path

	fmt.Println("PUT", url, payload)

	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("RESPONSE", resp.Status)

	return resp, nil
}

func (api *ApiClient) Delete(path string) (*http.Response, error) {
	url := api.baseUrl + path

	fmt.Println("DELETE", url)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("RESPONSE", resp.Status)

	return resp, nil
}

func (api *ApiClient) ParseBody(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("BODY %s", body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func assertStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		t.Fatalf(
			"Invalid Status Code. Expected Status \"%d\" and received \"%s\"",
			expected,
			resp.Status,
		)
	}
}
