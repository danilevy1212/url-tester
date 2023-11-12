package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type HeadResult struct {
	URL        string
	StatusCode int
}

func parseArgs() (string, error) {
	flag.Parse()

	filename := flag.Arg(0)

	if filename == "" {
		return "", fmt.Errorf("no filepath specified")
	}

	path, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("invalid filename %s: %s", filename, err)
	}

	fileinfo, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("file %s does not exist: %s", path, err)
	}
	if fileinfo.IsDir() {
		return "", fmt.Errorf("given path %s is a directory", path)
	}

	return path, nil
}

func parseURLFile(fileHandler *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(fileHandler)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	return lines, nil
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

func worker(data <-chan string, results chan HeadResult, wg *sync.WaitGroup) {
	for url := range data {
		searchWithHead(url, wg, results)
	}
}

func main() {
	path, err := parseArgs()
	if err != nil {
		log.Fatal("Error:", err)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	urls, err := parseURLFile(file)
	if err != nil {
		log.Fatal("Error parsing URL file:", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))
	numWorkers := runtime.NumCPU() * 2
	results := make(chan HeadResult, len(urls))
	url := make(chan string, len(urls))

	// Setup workers
	for i := 0; i < numWorkers; i++ {
		go worker(url, results, &wg)
	}

	// Send urls to them
	for _, u := range urls {
		url <- u
	}
	close(url)

	// When all urls are checked, it is safe to close.
	wg.Wait()
	close(results)

	for r := range results {
		fmt.Printf("Status '%d': %s\n", r.StatusCode, r.URL)
	}
}
