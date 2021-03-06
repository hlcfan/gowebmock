package webmock_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hlcfan/webmock"
)

func TestWebMock(t *testing.T) {
	server := webmock.New()
	baseURL := server.URL()
	fmt.Println("===", baseURL)
	server.Start()

	client := &http.Client{}

	t.Run("It serves stub http requests with GET", func(t *testing.T) {
		server.Stub("GET", "/abc", "ok")

		resp, err := http.Get(baseURL + "/abc")
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusOK, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody := string(respBodyBytes)
		if respBody != "ok" {
			t.Errorf("unexpected response body, want: %s, got: %s", "ok", respBody)
		}
	})

	t.Run("It serves stub http requests with POST", func(t *testing.T) {
		server.Stub("POST", "/post", "ok post")

		payload := []byte(`{"name":"Alex"}`)
		req, err := http.NewRequest("POST", baseURL+"/post", bytes.NewBuffer(payload))
		if err != nil {
			panic(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusOK, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody := string(respBodyBytes)
		if respBody != "ok post" {
			t.Errorf("unexpected response body, want: %s, got: %s", "ok", respBody)
		}
	})

	t.Run("It serves stub http requests when query parameters match", func(t *testing.T) {
		url := "/get?foo=bar&a=b"
		response := "ok with query parameters"

		server.Stub("GET", url, response)

		resp, err := http.Get(baseURL + "/get?foo=bar&a=b")
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusOK, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody := string(respBodyBytes)
		if respBody != response {
			t.Errorf("unexpected response body, want: %s, got: %s", "ok", respBody)
		}
	})

	t.Run("It serves stub http requests when query parameters don't match", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/get?foo=bar")
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("It serves stub http requests when headers match", func(t *testing.T) {
		url := "/get"
		response := "ok with headers"

		server.Stub("GET", url, response, webmock.WithHeaders("Accept-Encoding: gzip,deflate"))

		req, err := http.NewRequest("GET", baseURL+url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Accept-Encoding", "gzip,deflate")

		resp, err := client.Do(req)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusOK, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody := string(respBodyBytes)
		if respBody != response {
			t.Errorf("unexpected response body, want: %s, got: %s", "ok", respBody)
		}
	})

	t.Run("It serves stub http requests when headers don't match", func(t *testing.T) {
		url := "/get"
		response := "ok with headers"

		server.Stub("GET", url, response, webmock.WithHeaders("Accept-Encoding: gzip,deflate"))

		req, err := http.NewRequest("GET", baseURL+url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Accept-Encoding", "gzip")

		resp, err := client.Do(req)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("It serves stub http requests with customized response", func(t *testing.T) {
		url := "/get"
		response := "No permissions"

		server.Stub(
			"GET",
			url,
			"",
			webmock.WithResponse(http.StatusUnauthorized, response, map[string]string{
				"Access-Control-Allow-Origin": "*",
			}),
		)

		req, err := http.NewRequest("GET", baseURL+url, nil)
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusUnauthorized, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody := string(respBodyBytes)
		if respBody != response {
			t.Errorf("unexpected response body, want: %s, got: %s", "ok", respBody)
		}

		responseHeader := resp.Header.Get("Access-Control-Allow-Origin")
		if responseHeader != "*" {
			t.Errorf("unexpected response header, want: %s, got: %s", "*", responseHeader)
		}
	})

}

func TestWebMockLoadCassettes(t *testing.T) {
	server := webmock.New()
	baseURL := server.URL()
	fmt.Println("===", baseURL)
	server.Start()

	client := &http.Client{}

	t.Run("It serves stub http requests with cassette file", func(t *testing.T) {
		response := "OK, zoomer"
		// server.LoadCassettes("./fixtures/sample_cassette.yml")
		server.LoadCassettes("./fixtures")

		req, err := http.NewRequest("GET", baseURL+"/hello", nil)
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusOK, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody := string(respBodyBytes)
		if respBody != response {
			t.Errorf("unexpected response body, want: %s, got: %s", "ok", respBody)
		}

		responseHeader := resp.Header.Get("Access-Control-Allow-Origin")
		if responseHeader != "*" {
			t.Errorf("unexpected response header, want: %s, got: %s", "*", responseHeader)
		}

		requestID := "fake-request-id"
		requestIDHeader := resp.Header.Get("X-Request-Id")
		if requestIDHeader != requestID {
			t.Errorf("unexpected response header, want: %s, got: %s", requestID, requestIDHeader)
		}

		response = "Book created"

		req, err = http.NewRequest("POST", baseURL+"/book", nil)
		if err != nil {
			panic(err)
		}

		resp, err = client.Do(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusCreated, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody = string(respBodyBytes)
		if respBody != response {
			t.Errorf("unexpected response body, want: %s, got: %s", response, respBody)
		}

		response = "Service Unavailable"
		req, err = http.NewRequest("GET", baseURL+"/maintenance?foo=bar", nil)
		if err != nil {
			panic(err)
		}

		resp, err = client.Do(req)
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusServiceUnavailable, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody = string(respBodyBytes)
		if respBody != response {
			t.Errorf("unexpected response body, want: %s, got: %s", response, respBody)
		}

	})
}

func TestWebMockReset(t *testing.T) {
	server := webmock.New()
	baseURL := server.URL()
	fmt.Println("===", baseURL)
	server.Start()

	t.Run("It resets all routes", func(t *testing.T) {
		server.Stub("GET", "/abc", "ok")

		resp, err := http.Get(baseURL + "/abc")
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusOK, resp.StatusCode)
		}

		defer resp.Body.Close()

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		respBody := string(respBodyBytes)
		if respBody != "ok" {
			t.Errorf("unexpected response body, want: %s, got: %s", "ok", respBody)
		}

		server.Reset()

		resp, err = http.Get(baseURL + "/abc")
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("unexpected response status, want: %d, got: %d", http.StatusNotFound, resp.StatusCode)
		}
	})
}
