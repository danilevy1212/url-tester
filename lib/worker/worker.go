package worker

import (
	"context"
	"net/http"
	"sync"
	"time"
	"log"
)

type HeadResult struct {
	URL        string
	StatusCode int
}

func searchWithHead(url string, wg *sync.WaitGroup, results chan<- HeadResult) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		log.Printf("Error creating request: %s\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		results <- HeadResult{
			URL:        url,
			StatusCode: resp.StatusCode,
		}
	}
}

func Worker(data <-chan string, results chan HeadResult, wg *sync.WaitGroup) {
	for url := range data {
		searchWithHead(url, wg, results)
	}
}
