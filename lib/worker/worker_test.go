package worker

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

func TestWorkerHttpStatusOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	urls := make(chan string, 1)
	results := make(chan HeadResult, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go Worker(urls, results, &wg)

	urls <- server.URL
	close(urls)

	wg.Wait()
	close(results)

	_, ok := <-results
	if ok {
		t.Fatal("URL failed")
	}
}

func TestWorkerHttpStatusNotOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	urls := make(chan string, 1)
	results := make(chan HeadResult, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go Worker(urls, results, &wg)

	urls <- server.URL
	close(urls)

	wg.Wait()
	close(results)

	result, ok := <-results
	if !ok {
		t.Fatal("Expected a result, but didn't get one")
	}

	if result.URL != server.URL || result.StatusCode != http.StatusNotFound {
		t.Errorf("Expected URL: %s with StatusCode: %d; got URL: %s with StatusCode: %d", server.URL, http.StatusNotFound, result.URL, result.StatusCode)
	}
}

func TestWorkerHttpMultipleStatus(t *testing.T) {
	activePath, inactivePath := "/active", "/inactive"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == inactivePath {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.URL.Path == activePath {
			w.WriteHeader(http.StatusOK)
			return
		}

		t.Fatalf("Unimplemented path handler for %s", r.URL.Path)
	}))
	defer server.Close()

	activeURL, _ := url.JoinPath(server.URL, activePath)
	inactiveURL, _ := url.JoinPath(server.URL, inactivePath)

	testURLs := []string{
		activeURL,
		inactiveURL,
	}

	urls := make(chan string, len(testURLs))
	results := make(chan HeadResult, len(testURLs))
	var wg sync.WaitGroup

	wg.Add(len(testURLs))
	go Worker(urls, results, &wg)

	for _, testURL := range testURLs {
		urls <- testURL
	}
	close(urls)

	wg.Wait()
	close(results)

	if len(results) != 1 {
		t.Fatal("Only one URL should have not been found")
	}

	result, ok := <-results
	if !ok {
		t.Fatal("Expected a result, but didn't get one")
	}

	if result.URL != inactiveURL || result.StatusCode != http.StatusNotFound {
		t.Errorf("Expected URL: %s with StatusCode: %d; got URL: %s with StatusCode: %d", inactiveURL, http.StatusNotFound, result.URL, result.StatusCode)
	}
}
